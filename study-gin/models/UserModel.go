package model

import (
	"log"

	"gorm.io/gorm"
)

type TestModel struct {
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
}

func NewTestModel() *TestModel {
	return &TestModel{
		Name:    "",
		Age:     0,
		Address: "",
	}
}

type UserModel struct {
	gorm.Model
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	// omitempty可以将字段为空值的, nil的值忽略掉, 直接就不反了
	Role int `json:",omitempty"`
}
type User UserModel

func (u *User) GetUserList() interface{} {
	// 通过扩展user, 来增加一个反给前端的临时字段Address,(这个字段在数据库是不存在的)
	// 这个扩展只需将User作为tempUser的匿名字段即可, 方便后续处理
	type tempUser struct {
		User
		Address string `json:"address"`
		//gorm.DeletedAt `json:"del_time"`
	}
	// 查询单条数据
	_ = dbreader.Take(&u)

	// 查询全部数据
	var userList []User
	var resultList []interface{}
	dbreader.Find(&userList)

	for _, row := range userList {
		temp := new(tempUser)
		// 将原始的row赋给temp的User字段
		temp.User = row

		// 增加一个Address字段
		temp.Address = "河北省衡水市"

		/**
		这里一个思路, 结构体是不方便直接增加, 减少字段的,
		可以考虑将结构体变为map, 这样变动字段就会方便的多
		*/

		// 去掉row里的deleteAt字段
		//delete(temp.DeletedAt)

		// 假设不想反给前端这个deleteat字段, 则根据UserModel里定义的
		// tag `omitempty` 将deleteat置为空值
		// temp.Model.DeleteAt = ""
		resultList = append(resultList, temp)
	}

	return map[string]interface{}{
		"result_list": resultList,
		//"data_list": userList,
		//"userListType": reflect.TypeOf(userList),
	}
}

func (u *User) CreateUser() string {
	user := User{Username: "llf", Password: "123444", Role: 1}

	err := dbreader.AutoMigrate(&User{})

	if err != nil {
		log.Fatal("auto migrate err: ", err)
	}

	result1 := dbreader.Create(&user) // 通过数据的指针来创建

	log.Printf("insert user id %#v \n", user.ID)
	log.Printf("result is %#v \n", result1)
	log.Printf("result.error is %#v \n", result1.Error)

	// 用指定的字段创建记录

	result2 := dbreader.Select("Username", "Password", "Role").Create(&user)

	log.Printf("=====> user %#v\n", user)
	log.Printf("result2: %v \n", result2)
	log.Printf("result2 error %v \n", result2.Error)

	return "create user ok"
}

// d run -id --name mysql3308 \
// -p 3308
// -v /docker/mysql/slave_3308/mysql:/etc/mysql \
// -e MYSQL_ROOT_PASSWORD=root \
// --link mysql3306 \
// mysql:5.6

/**
1. 搭建好两套一模一样的mysql环境
2. 主从服务器分别开启binlog
		[
			1. log-bin=mysql-master01-bin       // 这里的这个值可以随意指定，最好有意义的master
			2. server-id=1                      // 这里的这个server-id也可以随意指定， 但是不能和其他服务器重复
		]
3. 配置需要同步的那些库
	binlog-do-db=blog  // 这个参数如果不指定， 默认同步全部的库
4. 配置屏蔽系统默认的数据库
	binlog-ignore-db=mysql
	binlog-ignore-db=information_schema
	binlog-ignore-db=performance_schema
5.  从服务和主服务器类似的配置， 以上3步再来一遍， // 这里有误，不必配置第4条也可以
6. 登录主服务器，创建一个用来同步数据的mysql账户
	（这里实际应用中， 还需要创建一个用户只具备read的权限， 用来在代码里连接mysql从库， 而这个db_sync的用户仅用来主从服务器之间同步）
    grant replication slave on *.* to 'DB_SYNC'@'%' identified by '123456';   // 这里的DB_SYNC是一个自定义的用户名， identitfied by 后面跟的是登录密码
	flush privileges;


7， 执行 show master status; 在查询结果中记录下来file 和position两个字段的值， 后续要用到， 该字段用来确认从服务器根据那个log文件来同步数据

8. 登录从服务器， 执行以下几步
	[
		1. stop slave // 先停止从服务器
		// 这里注意， 是整型的值就必须写整型，不能有引号包裹，例如：mysql_port=3306就不能写成mysql_port="3306",
		2. change master to master_host='master', master_user='DB_SYNC', master_password='123456', master_port=3306, master_log_file='mysql-bin.000014', master_log_pos=154, master_connect_retry=30;
		3. start slave
		4. show slave status // 检查从服务器是否启动了
	]


9. 给mysql创建用户：
	 create user 用户名 identified by '密码';
   精准控制用户的权限：
     grant 权限 on 数据库.数据表 to '用户' @ '主机名';
	 例如：
	 grant select on *.* to 'reader'@'%'；
     grant select,insert,update,delete,create,drop on vtdc.employee to joe@10.163.225.87 identified by ‘123′;

*/
