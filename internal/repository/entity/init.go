package entity

import "gorm.io/gorm"

// 自动初始化user表
func InitTable(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
