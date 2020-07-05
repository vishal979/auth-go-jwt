package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/models"
	"github.com/vishal979/auth/server/responses"
)

type password struct {
	ID       uint64 `json:"id"`
	Password string `json:"password"`
}

func (server *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	log.Info("user unmarshal successful", user)
	log.Info("preparing user for login")
	if err != nil {
		log.Error("error preparing user for login", err)
		responses.ERROR(w, 500, err)
		return
	}
	log.Info("user prepared to login")
	log.Info("validating user")
	err = user.Validate()
	if err != nil {
		log.Error("error while log in user", err)
		responses.ERROR(w, 500, err)
		return
	}
	log.Info("valdiating user successful")
	tokens, err := server.signIn(user.Username, user.Password)
	if err != nil {
		log.Error("error while signing in user", err)
		responses.ERROR(w, 500, err)
		return
	}
	responses.JSON(w, 200, tokens)
}
