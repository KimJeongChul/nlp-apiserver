package api

import (
	"net/http"
	"nlp-apiserver/config"
	"nlp-apiserver/errors"
	"nlp-apiserver/logger"
	"nlp-apiserver/utils"

	"github.com/gorilla/mux"
)

const (
	packageName = "api"
)

// ApiServer
type ApiServer struct {
	serverConfig *config.ServerConfigJson
	router       *mux.Router
}

var uuidCnt int
var localMacAddr string

func init() {
	uuidCnt = 0
	localMacAddrs, errMac := utils.GetMacAddress()
	if errMac != nil {
		logger.LogE(packageName, "init", "utils.GetMacAddress Err Msg=", errMac.Error())
		panic(errMac)
	}
	localMacAddr = localMacAddrs[0]
}

// NewApiServer Create new ApiServer
func NewApiServer(configJson *config.ServerConfigJson) *ApiServer {
	var as ApiServer
	as.serverConfig = configJson

	as.router = mux.NewRouter()
	as.router.HandleFunc("/train", as.handleTrainCreate)

	return &as
}

// Listen Create http, https server and ListenAndServe
func (as *ApiServer) Listen() int {
	funcName := "Listen"
	var err error
	if as.serverConfig.Ssl == 0 {
		logger.LogI(packageName, funcName, "HTTP Listening:Port=", as.serverConfig.ListenPort)
		err = http.ListenAndServe(":"+as.serverConfig.ListenPort, as.router)
		cErr := errors.NewCError(errors.HTTP_SERVE_ERR, err.Error())
		logger.LogE(packageName, funcName, "ERR:Msg=", cErr.Error())
	} else {
		logger.LogI(packageName, funcName, "HTTPs Listening:Cert=", as.serverConfig.CertPemPath, ",keyPem=", as.serverConfig.KeyPemPath)
		err = http.ListenAndServeTLS(":"+as.serverConfig.ListenPort, as.serverConfig.CertPemPath, as.serverConfig.KeyPemPath, as.router)
		cErr := errors.NewCError(errors.HTTPS_SERVE_ERR, err.Error())
		logger.LogE(packageName, funcName, "ERR:Msg=", cErr.Error())
	}
	if err != nil {
		panic(err)
	}
	return 0
}
