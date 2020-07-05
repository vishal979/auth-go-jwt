package controller

import (
	"errors"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/models"
	"github.com/vishal979/auth/server/auth"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) persistAuth(userid uint64, td *auth.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	err := server.RedisClient.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if err != nil {
		log.Info("error while setting value in redis")
	}
	err = server.RedisClient.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if err != nil {
		log.Info("error while setting value in redis")
	}
	return nil
}

func (server *Server) fetchAuth(authD *auth.AccessDetails) (uint64, error) {
	userid, err := server.RedisClient.Get(authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func (server *Server) deleteAuth(givenUUID string) (int64, error) {
	deleted, err := server.RedisClient.Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (server *Server) checkForExist(user models.User) error {
	log.Info("Checking if the email or username already exist")
	checkQuery := "SELECT * FROM Users WHERE username=?"
	stmt, err := server.DB.Prepare(checkQuery)
	if err != nil {
		log.Error("error preparing checking query ", err)
		if stmt != nil {
			defer stmt.Close()
		}
		return err
	}
	defer stmt.Close()
	result, err := stmt.Query(user.Username)
	if err != nil {
		log.Error("error executing checking query ", err)
		return err
	}
	if result.Next() {
		log.Error("email or username or phone already exists")
		return errors.New("email or username or phone already exists")
	}
	return nil
}

func (server *Server) createEntry(user models.User) error {
	insertQuery := "Insert into Users (username,password,) VALUES(?,?)"
	stmt, err := server.DB.Prepare(insertQuery)
	if err != nil {
		log.Error("error while preparing statement while creating entry", err)
		if stmt != nil {
			defer stmt.Close()
		}
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.Username, user.Password)
	if err != nil {
		log.Error("error inserting user ", err)
		return err
	}
	log.Info(result)
	return nil
}

func (server *Server) verifyPassword(username, pass string) (uint64, error) {
	var password password
	passwordQuery := "select id, password from Users where username=?"
	stmt, err := server.DB.Prepare(passwordQuery)
	if err != nil {
		log.Error("error preparing checking query", err)
		if stmt != nil {
			defer stmt.Close()
		}
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Query(username)
	log.Info("email", username)
	if err != nil {
		log.Error("error executing checking query", err)
		return 0, err
	}
	defer result.Close()
	if !result.Next() {
		log.Error("Unable to find the user with email", username)
		return 0, errors.New("Unable to find the user")
	}
	result.Scan(&password.ID, &password.Password)
	err = bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(pass))
	if err != nil {
		log.Error("password mismatch error")
		return 0, errors.New("password mismatch")
	}
	return password.ID, nil
}

func (server *Server) signIn(username, password string) (map[string]string, error) {
	var err error
	log.Info("trying to sign in user")
	id, err := server.verifyPassword(username, password)
	if err != nil {
		return nil, err
	}
	log.Info("password verification successful")
	log.Info("creating token")
	ts, err := auth.CreateToken(id)
	if err != nil {
		log.Info("error creating token")
		return nil, err
	}
	err = server.persistAuth(id, ts)
	if err != nil {
		return nil, err
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return tokens, nil

}
