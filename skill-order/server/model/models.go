package model

// 定义商品表模型
type SkillGoods struct {
	Id        int     `gorm:"primarykey"`
	Name      string  `gorm:"column:name;size:200" json:"name"` // 数据库字段名为 name，JSON 字段名为 name
	Price     float32 `gorm:"column:price;comment:单价;type:DECIMAL(10,2);not_null" json:"price"`
	Stock     int16   `gorm:"column:stock;comment:库存" json:"stock"`
	CreatedAt int64   `gorm:"autoCreateTime:milli"`
	UpdatedAt int64   `gorm:"index;autoUpdateTime:milli"`
}

// 定义订单模型
type Order struct {
	Id          int    `gorm:"primarykey"`
	OrderSn     string `gorm:"type:char(100);comment:订单号"`
	Uid         string `gorm:"type:char(100);comment:用户id"`
	GoodsId     int
	GoodsNumber int
	CreatedAt   int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64 `gorm:"index;autoUpdateTime:milli"`
}
