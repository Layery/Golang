package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

type Salary struct {
	totalSalary     int
	totalDays       int
	missHours       int
	forgetSignTimes int
}

func NewSalary(totalSalary, totalDays, missHours, forgetSignTimes int) *Salary {
	return &Salary{
		totalSalary: totalSalary,
		totalDays:   totalDays,
		missHours:    missHours,
		forgetSignTimes:    forgetSignTimes,
	}
}

func (s Salary) GetSalary() error {
	mylog(fmt.Sprintf("本月实际出勤%d天", s.totalDays))
	missSalary := 0
	miss4Sign := 0
	if s.totalDays <= 0 {
		return errors.New("totalDays不能为0!")
	}
	if s.missHours > 0 {
		mylog(fmt.Sprintf("本月请假共%d小时", s.missHours))
		missSalary = (s.totalSalary/s.totalDays/8) * s.missHours
	}
	if s.forgetSignTimes > 0 {
		mylog(fmt.Sprintf("本月%d次忘打卡", s.forgetSignTimes))
		miss4Sign = s.forgetSignTimes*10
	}

	rs := s.totalSalary - missSalary - miss4Sign
	mylog(rs)
	return nil
}

func mylog(data interface{})  {
	// 这个格式化的模板是真特么奇葩
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%v %v \n", t, data)
}

func main()  {
	app := cli.App{
		Name:                   "SalaryCalculator",
		Usage:                  "SalaryCalculator used to calculate you salary",
		Version:                "1.0.0",
		Description:            "因公司奇葩的工资计算方式, SalaryCalculator应运而生",
		Commands:               nil,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "salary",
				Aliases:     []string{"s"},
				Value: 5000,
				Usage:       "`salary` 你和老板约定的薪资",
				Required:    true,
			},
			&cli.IntFlag{
				Name:    "days",
				Aliases: []string{"d"},
				Usage:   "`days` 当月实际上班的天数",
				Required:    true,
			},
			&cli.IntFlag{
				Name:    "missHours",
				Aliases: []string{"m"},
				Value:   0,
				Usage:   "`missHours` 当月请假的时长, 单位为: 小时, 请假一天即8小时",
			},
			&cli.IntFlag{
				Name:    "forgetSignTimes",
				Aliases: []string{"f"},
				Value:   0,
				Usage:   "`forgetSign` 当月忘打卡的次数",
			},
		},
		Action: func(context *cli.Context) error {
			totalSalary := context.Int("salary")
			totalDays := context.Int("days")
			missHours := context.Int("missHours")
			forgetSignTimes := context.Int("forgetSignTimes")
			return NewSalary(totalSalary, totalDays, missHours, forgetSignTimes).GetSalary()
		},
		Compiled: time.Time{},
		Authors: []*cli.Author{
			&cli.Author{Name: "weidingyi", Email: "weidingyi@qq.com"},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}