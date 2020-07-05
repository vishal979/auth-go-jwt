package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/models"
	"github.com/vishal979/auth/server/responses"
)

func (server *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Info("preparing user for signup")
	err = user.Prepare("signup")
	if err != nil {
		log.Error("error preparing user for signup", err)
		responses.ERROR(w, 500, err)
		return
	}
	log.Info("user prepared to signup")
	log.Info("validating user")
	err = user.Validate()
	if err != nil {
		log.Error("Error validating user for signup", err)
		responses.ERROR(w, 500, err)
		return
	}
	log.Info("User validation successful")
	if err := server.checkForExist(user); err != nil {
		responses.ERROR(w, 500, errors.New("duplicate record exists"))
		return
	}
	log.Info("checking user for duplicate record successful, no duplicate records for the user")
	//create an entry for the user
	if err := server.createEntry(user); err != nil {
		responses.ERROR(w, 500, errors.New("error while creating user"))
	}
	responses.JSON(w, 200, "User Created Successfully")
}
