package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/server/auth"
	"github.com/vishal979/auth/server/responses"
)

func (server *Server) refreshtokenHandler(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	log.Info("Post Request refresh token")
	log.Info(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &mapToken)
	if err != nil {
		log.Error("refresh unmarshal error")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	refreshToken := mapToken["refresh_token"]
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		log.Error("refresh token expired")
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		log.Error("unauthorized user")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			log.Error("error while refresh uuid")
			return
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			log.Error("userid refreshtoken.go error")
			return
		}
		deleted, delErr := server.deleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			log.Error("unauthorized user")
			return
		}
		ts, createErr := auth.CreateToken(userID)
		if createErr != nil {
			log.Error("create error refreshtoken.go")
			return
		}
		saveErr := server.persistAuth(userID, ts)
		if saveErr != nil {
			log.Error("error saving refreshtoken.go")
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		log.Info(http.StatusCreated, tokens)
	} else {
		log.Error("refresh expired")
	}
}
