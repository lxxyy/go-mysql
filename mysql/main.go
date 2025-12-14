package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // 导入数据库驱动
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func initDB() (err error) {
	// 初始化数据库连接
	dsn := "root:Liu12345@tcp(127.0.0.1:3306)/go_test"
	// 打开数据库 再次强调 这里不能使用 := 赋值 因为是全局变量
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		//log.Fatal("数据库连接失败:", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		//log.Fatal("数据库连接失败:", err)
		return err
	}
	db.SetMaxOpenConns(10)
	// db.Query()
	//fmt.Println("数据库连接成功！")
	return
}

// queryUserOne 查询单个用户
func queryUserOne(id int) (user *User, err error) {
	user = &User{}
	sqlStr := `select id,name,age from user where id=?;`
	err = db.QueryRow(sqlStr, id).Scan(&user.Id, &user.Name, &user.Age)
	if err != nil {
		fmt.Println("查询失败:", err)
		return nil, err
	}

	fmt.Printf("id:%d name:%s age:%d\n", user.Id, user.Name, user.Age)
	return user, nil
}

// queryByIdList 查询id列表
func queryByIdList(id int) (users []*User, err error) {
	users = []*User{}
	sqlStr := `select id,name,age from user where id > ?;`
	rows, err := db.Query(sqlStr, id)
	if err != nil {
		fmt.Println("查询失败:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			fmt.Println("scan failed:", err)
			return nil, err
		}
		users = append(users, user)
		fmt.Printf("id:%d name:%s age:%d\n", user.Id, user.Name, user.Age)
	}
	return users, nil
}

// intertUser 插入数据
func intertUser(name string, age int) {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, name, age)
	if err != nil {
		fmt.Println("insert failed:", err)
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Println("get lastinsert ID failed:", err)
	}
	fmt.Println("insert success, id:", id)
}

// updateUser 更新
func updateUser(id int, age int) {
	sqlStr := "update user set age=? where id=?"
	ret, err := db.Exec(sqlStr, age, id)
	if err != nil {
		fmt.Println("update failed:", err)
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected failed:", err)
	}
	fmt.Println("update success, affected rows:", n)
}

// deleteUser 删除
func deleteUser(id int) {
	sqlStr := "delete from user where id=?"
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		fmt.Println("delete failed:", err)
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected failed:", err)
	}
	fmt.Println("delete success, affected rows:", n)
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功！")
	user, err := queryUserOne(1)
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	fmt.Println("查询成功:", *user)

	userList, err := queryByIdList(0)
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	fmt.Println("查询成功:", userList)

	//intertUser("王五", 22)
	//intertUser("赵六", 23)
	// updateUser(1, 25)
	// deleteUser(5)
}
