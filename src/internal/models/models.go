package models

import (
	"time"

	"github.com/google/uuid"
)

type DocumentMetaData struct {
	ID           uuid.UUID
	PageCount    int
	DocumentName string
	CreatorID    uint64
	CreationTime time.Time
}

type DocumentData struct {
	ID            uuid.UUID
	DocumentBytes []byte //pdf file -- the whole file
}

type Markup struct {
	ID         uint64
	PageData   []byte    //png file -- the page data
	ErrorBB    []float32 //Bounding boxes in yolov8 format
	ClassLabel uint64
	CreatorID  uint64
}

type Role int

const (
	Sender Role = iota // Role check depends on the order
	Controller
	Admin
)

func (r Role) ToString() string {

	switch r {
	case Sender:
		return "Sender"
	case Controller:
		return "Controller"
	case Admin:
		return "Admin"

	default:
		return "Unknown"
	}

}

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
	ID          uint64
	ClassName   string
}

type ErrorReport struct {
	DocumentID uuid.UUID
	ReportData []byte
}

type Token struct {
	UserID  uint64
	ExpTime time.Duration // think about securing cookies, store cookies on backend (hashing or storing)
	Role    Role
}
