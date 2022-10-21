package pojo

import "Golangapi/config"

type User struct {
	First_name      string `json:"first_name" binding:"required"`
	Last_name       string `json:"last_name" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Username        string `json:"username" binding:"required,email"`
	ID              string `json:"id"`
	Account_created string `json:"account_created"`
	Account_updated string `json:"account_updated"`
}

func GetAllUsers_db() []User {
	var users []User
	config.DB.Find(&users)
	return users
}

func GetUsers_db(userID string) User {
	var user User
	config.DB.Where("id=?", userID).First(&user)
	return user
}

func GetUsers_db_Username(userID string) string {
	var user User
	config.DB.Where("id=?", userID).First(&user)
	return user.Username
}

func GetUsers_db_Pass(userID string) string {
	var user User
	config.DB.Where("id=?", userID).First(&user)
	return user.Password
}

func PostUsers_db(user User) User {
	config.DB.Create(&user)
	return user
}

func DeleteUser(userID string) {
	user := User{}
	config.DB.Where("id=?", userID).Delete(&user)
}

func UpdateUser(userID string, user User) User {
	config.DB.Model(&user).Where("id=?", userID).Updates(user)
	return user
}
