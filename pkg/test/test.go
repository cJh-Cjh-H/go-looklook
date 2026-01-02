package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:Mcjh0908.@tcp(139.9.177.213:3306)/LookLook?charset=utf8mb4&parseTime=true&loc=Local"

	// 测试基本连接
	fmt.Println("🔍 开始数据库诊断...")

	// 1. 测试连接
	fmt.Println("\n1. 测试连接...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping失败:", err)
	}
	fmt.Println("✅ 连接成功")

	// 2. 测试查询性能
	fmt.Println("\n2. 测试查询性能...")
	start := time.Now()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables").Scan(&count)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("✅ 查询完成，耗时: %v，结果: %d\n", time.Since(start), count)

	// 3. 检查连接池
	fmt.Println("\n3. 检查连接池状态...")
	stats := db.Stats()
	fmt.Printf("打开连接数: %d\n", stats.OpenConnections)
	fmt.Printf("使用中连接: %d\n", stats.InUse)
	fmt.Printf("空闲连接: %d\n", stats.Idle)

	// 4. 模拟并发查询
	fmt.Println("\n4. 模拟并发查询...")
	testConcurrentQueries(db)

	fmt.Println("\n🎉 诊断完成")
}

func testConcurrentQueries(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		go func(id int) {
			var result int
			err := db.QueryRowContext(ctx, "SELECT ?", id).Scan(&result)
			if err != nil {
				fmt.Printf("协程 %d 失败: %v\n", id, err)
			} else {
				fmt.Printf("协程 %d 成功\n", id)
			}
		}(i)
	}

	time.Sleep(2 * time.Second)
}
