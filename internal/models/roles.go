package models

type RoleStr string

const (
	RoleUser  RoleStr = "user"
	RoleAdmin RoleStr = "admin"
	RoleMod   RoleStr = "moderator"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Level       int    `json:"level"`
}
