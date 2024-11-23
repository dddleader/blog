package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultUsername = "admin"
	dbConnString    = "root:200455@tcp(127.0.0.1:3307)/blog"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run create_admin.go <password>")
	}

	password := os.Args[1]

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("生成密码哈希失败: %v", err)
	}

	// 连接数据库
	db, err := sql.Open("mysql", dbConnString)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 检查连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 插入管理员账户
	_, err = db.Exec(
		"INSERT INTO admin_users (username, password) VALUES (?, ?)",
		defaultUsername,
		string(hashedPassword),
	)
	if err != nil {
		log.Fatalf("插入管理员账户失败: %v", err)
	}

	fmt.Println("管理员账户创建成功!")
	fmt.Printf("用户名: %s\n", defaultUsername)
	fmt.Printf("密码哈希: %s\n", string(hashedPassword))
}
