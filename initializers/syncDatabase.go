package initializers

import "github.com/Ryuuukin/ap-assignment1/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate((&models.Post{}))
}
