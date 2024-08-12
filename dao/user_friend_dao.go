package dao

import "github.com/google/uuid"

type UserFriendDao struct {
	Id        uuid.UUID `json:"id" binding:"required,uuid"`
	Name      string    `json:"name" binding:"required"`
	Highscore int       `json:"highscore" binding:"required"`
}
