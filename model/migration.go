package model

import "log"

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
	}
}
