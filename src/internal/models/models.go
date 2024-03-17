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
type User struct {
	Login    string
	Password string
	Name     string
	Surname  string
	Role     string
	Group    string // in case it is a controller it will have work group, in case of user, his group
}

type MarkupType struct {
	Description string
	CreatorID   int
}
