package models

import (
	"time"
)

type Document struct {
	PageCount    int
	DocumentData []byte
	ChecksCount  int
	CreatorID    uint64
	CreationTime time.Time
}

type Markup struct {
	ErrorBB    []float32
	ClassLabel uint64
}

type Role int

const (
	user Role = iota
	admin
	controller
)

type User struct {
	ID       uint64
	Login    string
	Password string
	Name     string
	Surname  string
	Role     Role
	Group    string // in case it is a controller it will have work group, in case of user, his group
}

type MarkupType struct {
	Description string
	CreatorID   int
}

type Cookie struct {
	Token   string
	UserID  uint64
	ExpTime time.Duration
	Role    Role
}
