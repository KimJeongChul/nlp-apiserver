package api

import (
	"net/http"
	"nlp-apiserver/config"
	"nlp-apiserver/errors"
	"nlp-apiserver/logger"

	"github.com/gorilla/mux"
)

// ApiServer
type ApiServer struct {
	serverConfig *config.ServerConfigJson
	router       *mux.Router
}

// NewApiServer Create new ApiServer
func NewApiServer(configJson *config.ServerConfigJson) *ApiServer {
	var as ApiServer
	as.serverConfig = configJson

	as.router = mux.NewRouter()

	return &as
}

// Listen Create http, https server and ListenAndServe
func (as *ApiServer) Listen() int {
	var err error
	if as.serverConfig.Ssl == 0 {
		logger.LogI("Listen", "api", "HTTP Listening:Port=", as.serverConfig.ListenPort)
		err = http.ListenAndServe(":"+as.serverConfig.ListenPort, as.router)
		cErr := errors.NewCError(errors.HTTP_SERVE_ERR, err.Error())
		logger.LogE("Listen", "api", "ERR:Msg=", cErr.Error())
	} else {
		logger.LogI("Listen", "api", "HTTPs Listening:Cert=", as.serverConfig.CertPemPath, ",keyPem=", as.serverConfig.KeyPemPath)
		err = http.ListenAndServeTLS(":"+as.serverConfig.ListenPort, as.serverConfig.CertPemPath, as.serverConfig.KeyPemPath, as.router)
		cErr := errors.NewCError(errors.HTTPS_SERVE_ERR, err.Error())
		logger.LogE("Listen", "api", "ERR:Msg=", cErr.Error())
	}
	if err != nil {
		panic(err)
	}
	return 0
}
