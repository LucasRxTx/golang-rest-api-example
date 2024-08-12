package dao

import (
	"github.com/google/uuid"
)

type UserDao struct {
	Id   uuid.UUID
	Name string
}
