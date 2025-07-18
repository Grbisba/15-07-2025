package dto

import (
	"github.com/google/uuid"
)

type (
	CreateUser struct {
		Login    string
		Password string
	}
	Login struct {
		Login    string
		Password string
	}
	User struct {
		ID          uuid.UUID
		Login       string
		EncPassword string
		TotalTasks  int
	}
)
