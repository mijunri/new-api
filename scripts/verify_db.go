package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func verifyDatabase(host, port, user, password, dbname string) error {
	// 构建 MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	fmt.Printf("正在连接数据库: %s@%s:%s/%s\n", user, host, port, dbname)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	defer sqlDB.Close()

	// 查询所有表
	var tables []string
	err = db.Raw("SHOW TABLES").Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("查询表列表失败: %v", err)
	}

	fmt.Printf("\n数据库 %s 中的表:\n", dbname)
	fmt.Println("=" + string(make([]byte, 40)))

	expectedTables := []string{
		"abilities",
		"channels",
		"checkins",
		"logs",
		"midjourneys",
		"models",
		"options",
		"passkey_credentials",
		"prefill_groups",
		"quota_data",
		"redemptions",
		"setups",
		"tasks",
		"tokens",
		"top_ups",
		"two_fa_backup_codes",
		"two_fas",
		"users",
		"vendors",
	}

	// 创建表名映射
	tableMap := make(map[string]bool)
	for _, t := range tables {
		tableMap[t] = true
	}

	fmt.Println("\n已创建的表:")
	for _, t := range tables {
		fmt.Printf("  ✓ %s\n", t)
	}

	fmt.Println("\n表创建状态检查:")
	missingTables := []string{}
	for _, expected := range expectedTables {
		if tableMap[expected] {
			fmt.Printf("  ✓ %s - 存在\n", expected)
		} else {
			fmt.Printf("  ✗ %s - 缺失\n", expected)
			missingTables = append(missingTables, expected)
		}
	}

	if len(missingTables) > 0 {
		fmt.Printf("\n警告: 以下表缺失: %v\n", missingTables)
	} else {
		fmt.Println("\n所有必需的表都已创建!")
	}

	// 统计每个表的行数
	fmt.Println("\n表行数统计:")
	for _, t := range tables {
		var count int64
		err = db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM `%s`", t)).Scan(&count).Error
		if err != nil {
			fmt.Printf("  %s: 查询失败 - %v\n", t, err)
		} else {
			fmt.Printf("  %s: %d 行\n", t, count)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run verify_db.go [test|prod|all]")
		os.Exit(1)
	}

	env := os.Args[1]

	// 数据库配置
	host := "rm-t4n395l1o71a3h36vzo.mysql.singapore.rds.aliyuncs.com"
	port := "3306"
	user := "bitmodel"
	password := "bitmodel@1234"

	switch env {
	case "test":
		fmt.Println("=== 验证测试环境数据库 ===")
		err := verifyDatabase(host, port, user, password, "bitmodel")
		if err != nil {
			log.Fatalf("验证失败: %v", err)
		}

	case "prod":
		fmt.Println("=== 验证生产环境数据库 ===")
		err := verifyDatabase(host, port, user, password, "bitmodel-prod")
		if err != nil {
			log.Fatalf("验证失败: %v", err)
		}

	case "all":
		fmt.Println("=== 验证测试环境数据库 ===")
		err := verifyDatabase(host, port, user, password, "bitmodel")
		if err != nil {
			log.Printf("测试环境验证失败: %v", err)
		}
		fmt.Println("\n" + string(make([]byte, 50)) + "\n")

		fmt.Println("=== 验证生产环境数据库 ===")
		err = verifyDatabase(host, port, user, password, "bitmodel-prod")
		if err != nil {
			log.Printf("生产环境验证失败: %v", err)
		}

	default:
		fmt.Printf("未知的环境参数: %s\n", env)
		os.Exit(1)
	}

	fmt.Println("\n验证完成!")
}
