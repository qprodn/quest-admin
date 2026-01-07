package user

import "time"

type UpdatePasswordBO struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type UpdateStatusBO struct {
	UserID string
	Status int32
}

type UpdateLoginInfoBO struct {
	UserID    string
	LoginIP   string
	LoginDate time.Time
}

type ManageUserPostsBO struct {
	UserID    string
	PostIDs   []string
	Operation string
}
