package main

import (
	"flag"
	"nlp-apiserver/api"
	"nlp-apiserver/logger"
	"nlp-apiserver/utils"
	"os"
)

const (
	packageName = "main"
	funcName = "main"
)

func main() {
	// Load serverConfig.json
	serverConfigFilePath := flag.String("c", "./config/serverConfig.json", "Set serverConfig.json file path")
	serverConfig, errLoad := utils.LoadConfigJson(serverConfigFilePath)
	if errLoad != nil {
		logger.LogE(packageName, funcName, "Config file:" +  *serverConfigFilePath + " load error")
		os.Exit(-1)
	}

	if resultEnv := utils.GetEnv(serverConfig); resultEnv < 0 {
		os.Exit(-1)
	}

	logger.LogI(packageName, funcName, "Listening Config:HttpPort=", serverConfig.ListenPort)
	logger.LogI(packageName, funcName, "Listening Config:SSL=", serverConfig.Ssl)
	if serverConfig.Ssl == 1 {
		logger.LogI(packageName, funcName, "Listening Config:Cert=", serverConfig.CertPemPath, " Key=", serverConfig.KeyPemPath)
	}

	apiServer := api.NewApiServer(serverConfig)
	apiServer.Listen()
}