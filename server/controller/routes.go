package controller

import "github.com/vishal979/auth/server/middlewares"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.loginHandler)).Methods("POST")
	server.Router.HandleFunc("/logout", middlewares.SetMiddlewareAuthentication(server.logoutHandler)).Methods("POST")
	server.Router.HandleFunc("/signup", middlewares.SetMiddlewareJSON(server.signupHandler)).Methods("POST")
	server.Router.HandleFunc("/refreshtoken", middlewares.SetMiddlewareAuthentication(server.refreshtokenHandler)).Methods("POST")
	server.Router.HandleFunc("/validatetoken", middlewares.SetMiddlewareJSON(server.validateTokenHandler)).Methods("POST")
}
