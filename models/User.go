package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User struct for user
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Validate validates user
func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("empty fields")
	}
	return nil
}

//Prepare user
func (u *User) Prepare(action string) error {
	switch strings.ToLower(action) {
	case "login":
		u.ID = 0
		u.Username = html.EscapeString(strings.TrimSpace(u.Username))
		if u.Password == "" {
			return errors.New("password empty")
		}
		// hashedPassword, err := Hash(u.Password)
		// if err != nil {
		// 	log.Error("error generating hash", err)
		// 	return err
		// }
		// u.Password = string(hashedPassword)
		return nil
	case "signup":
		u.ID = 0
		u.Username = html.EscapeString(strings.TrimSpace(u.Username))
		hashedPassword, err := hash(u.Password)
		if err != nil {
			return fmt.Errorf("unable to calculate hash %e", err)
		}
		u.Password = string(hashedPassword)
		return nil
	default:
		return nil
	}
}

func hash(password string) ([]byte, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// VerifyPassword verifies hash to string
func VerifyPassword(hashedPassword, password string) error {
	log.Info([]byte(hashedPassword), "   ", []byte(password))
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
