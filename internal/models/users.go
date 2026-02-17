package models

import (
	"database/sql/driver"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type password struct {
	Hash []byte
	Text *string
}

func (p *password) Scan(src any) error {
	if src == nil {
		p.Hash = nil
		return nil
	}

	switch value := src.(type) {
	case []byte:
		p.Hash = append(p.Hash[:0], value...)
	case string:
		p.Hash = append(p.Hash[:0], []byte(value)...)
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *models.password", src)
	}

	return nil
}

func (p password) Value() (driver.Value, error) {
	if len(p.Hash) == 0 {
		return nil, nil
	}

	return p.Hash, nil
}

// Set hashes the plain password and stores it in the password struct
func (p *password) Set(plainPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	p.Text = &plainPassword
	return nil
}

// Check compares the provided plain password with the stored hash
func (p *password) Check(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainPassword))
	return err == nil
}
