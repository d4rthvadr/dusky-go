package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type password struct {
	hash []byte
	text *string
}

// Set hashes the plain password and stores it in the password struct
func (p *password) Set(plainPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	p.text = &plainPassword
	return nil
}

// Check compares the provided plain password with the stored hash
func (p *password) Check(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))
	return err == nil
}
