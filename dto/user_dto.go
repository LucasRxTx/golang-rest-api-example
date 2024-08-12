package dto

type UserDto struct {
	Id   string `json:"id" binding:"required,uuid"`
	Name string `json:"name" binding:"required"`
}
