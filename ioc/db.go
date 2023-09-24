package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/config"
	"webook/internal/repository/entity"
)

// 初始化数据库连接
func InitDB() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root:root@tcp(127.0.0.1:13306)/webook"
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		//	应该打日志
		// 初始数据库连接失败，server也没必要运行了
		panic(err)
	}

	// 自动迁移表结构
	err = entity.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
