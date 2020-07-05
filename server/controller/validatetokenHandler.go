package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/server/auth"
	"github.com/vishal979/auth/server/responses"
)

func (server *Server) validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		log.Println("error while extracting meta data from token")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	_, err = server.fetchAuth(tokenAuth)
	if err != nil {
		log.Error("error fetching")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
}
