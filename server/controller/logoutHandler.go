package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/server/auth"
)

func (server *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r)
	au, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		log.Error("unauthorized")
		return
	}
	deleted, delErr := server.deleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		log.Error("unauthorized")
		return
	}
	log.Error("logout successful")
}
