package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" // 导入数据库驱动
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sqlx.DB

func initDB() (err error) {
	// 初始化数据库连接
	dsn := "root:Liu12345@tcp(127.0.0.1:3306)/go_test"
	// 打开数据库 再次强调 这里不能使用 := 赋值 因为是全局变量
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//log.Fatal("数据库连接失败:", err)
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	// db.Query()
	//fmt.Println("数据库连接成功！")
	return
}
func queryUserObjOne(id int) (user *User, err error) {
	user = &User{}
	sqlStr := `select * from user where id=?;`
	err = db.Get(user, sqlStr, id)
	if err != nil {
		fmt.Println("查询失败:", err)
		return nil, err
	}

	fmt.Printf("id:%d name:%s age:%d\n", user.Id, user.Name, user.Age)
	return user, nil
}

func queryByIdList(id int) (users []*User, err error) {
	users = make([]*User, 10)
	sqlStr := `select * from user where id > ?;`
	err = db.Select(&users, sqlStr, id)
	if err != nil {
		fmt.Println("查询失败:", err)
		return nil, err
	}
	fmt.Printf("查询成功%#v\n", users)
	return users, nil
}
func main() {
	initDB()
	// queryUserObjOne(1)
	userList, err := queryByIdList(0)
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	fmt.Printf("查询成功%#v\n", userList)
	for _, user := range userList {
		fmt.Printf("查询成功%#v\n", *user)
	}
}
