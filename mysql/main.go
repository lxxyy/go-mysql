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

func prepareInsert(users []User) error {
	sqlStr := "insert into user(name, age) values(?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare failed:", err)
		return err
	}
	defer stmt.Close()

	for _, u := range users {
		ret, err := stmt.Exec(u.Name, u.Age)
		if err != nil {
			fmt.Println("insert failed:", err)
			return err
		}
		n, err := ret.RowsAffected()
		if err != nil {
			fmt.Println("get rows affected failed:", err)
			return err
		}
		fmt.Println("insert success, affected rows:", n)
	}
	//for k, v := range m {
	//
	//}

	return nil
}

// transactionExample 事务回滚案例
func transactionExample() error {
	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	// 使用defer确保在函数退出时根据情况提交或回滚事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	// 执行第一个操作 - 插入用户
	sqlStr1 := "INSERT INTO user(name, age) VALUES (?, ?)"
	_, err = tx.Exec(sqlStr1, "张三", 25)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("插入用户失败，事务已回滚: %v", err)
	}

	// 执行第二个操作 - 更新用户年龄
	sqlStr2 := "UPDATE usexxxr SET age = ? WHERE name = ?"
	_, err = tx.Exec(sqlStr2, 30, "李四")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新用户失败，事务已回滚: %v", err)
	}

	// 模拟一个错误情况触发回滚
	// 如果这里发生错误，事务会回滚
	simulatedError := true
	if simulatedError {
		tx.Rollback()
		return fmt.Errorf("模拟错误，事务已回滚")
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Println("事务执行成功")
	return nil
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功！")
	//user, err := queryUserOne(1)
	//if err != nil {
	//	fmt.Println("查询失败:", err)
	//}
	//fmt.Println("查询成功:", *user)
	//
	//userList, err := queryByIdList(0)
	//if err != nil {
	//	fmt.Println("查询失败:", err)
	//}
	//fmt.Println("查询成功:", userList)

	//intertUser("王五", 22)
	//intertUser("赵六", 23)
	// updateUser(1, 25)
	// deleteUser(5)

	//users := []User{
	//	{Name: "五千", Age: 25},
	//	{Name: "路费", Age: 30},
	//}
	//err = prepareInsert(users)
	//if err != nil {
	//	fmt.Println("批量插入失败:", err)
	//}

	transactionExample()
}
