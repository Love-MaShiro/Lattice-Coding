package main

import (
	"fmt"
	"log"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/db"
	"lattice-coding/internal/modules/provider/infra/persistence"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接 MySQL
	mysqlDB, err := db.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	fmt.Println("Connected to MySQL successfully")

	// 执行迁移
	fmt.Println("Starting database migration...")
	if err := persistence.Migrate(mysqlDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database migration completed successfully!")
	fmt.Println("\nTables created:")
	fmt.Println("1. providers - Provider 主表")
	fmt.Println("2. provider_healths - 健康状态独立表")
	fmt.Println("3. model_configs - 模型配置表")
}
