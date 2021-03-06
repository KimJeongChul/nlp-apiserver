package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"nlp-apiserver/config"
	"nlp-apiserver/logger"
	"os"
	"strconv"
	"time"
)

const (
	packageName = "utils"
)

// LoadConfigJson Load serverConfig.json
func LoadConfigJson(configFileName *string) (*config.ServerConfigJson, error) {
	funcName := "LoadConfigJson"
	var serverConfigJson config.ServerConfigJson
	file, errOpen := os.Open(*configFileName)
	if errOpen != nil {
		logger.LogE(packageName, funcName, "os.Open Path:", *configFileName, " Err Msg=", errOpen)
		return nil, errOpen
	}
	decoder := json.NewDecoder(file)
	errDecode := decoder.Decode(&serverConfigJson)
	if errDecode != nil {
		logger.LogE(packageName, funcName, "json.Decoder.Decode Err Msg=", errDecode)
		return nil, errDecode
	}

	return &serverConfigJson, nil
}

// GetEnv Get OS environments 
func GetEnv(config *config.ServerConfigJson) int {
	funcName := "GetEnv"
	// LISTEN_PORT
	if listenPort, isExist := osGetEnv("LISTEN_PORT"); isExist != true {
		logger.LogE(packageName, funcName, "os.GetEnv LISTEN_PORT value doesn't exist")
		return -1
	} else {
		config.ListenPort = listenPort
	}
	// SSL
	if ssl, isExist := osGetEnv("SSL"); isExist != true {
		logger.LogE(packageName, funcName, "os.GetEnv SSL value doesn't exist")
		return -1
	} else {
		issl, _ := strconv.Atoi(ssl)
		config.Ssl = issl
	}
	// LOG_LEVEL
	if level, isExist := osGetEnv("LOG_LEVEL"); isExist != true {
		logger.LogE(packageName, funcName, "os.GetEnv LOG_LEVEL value doesn't exist")
		return -1
	} else {
		config.LogLevel = level
	}
	return 0
}

// osGetEnv Get OS environments using key
func osGetEnv(key string) (string, bool) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", false
	}
	return value, true
}

// EnableCors
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

// GetMillisTimeFormat YYYYMMDDhhmmsslll
func GetMillisTimeFormat(t time.Time) string {
	// Golang timeformat 2006-01-02 15:04:05, Mon Jan 2 15:04:05 -0700 MST 2006
	timestamp := t.Format("20060102150405")
	return timestamp + strconv.Itoa(t.Nanosecond()/1000000)
}

// GetMacAddress Get MAC address
func GetMacAddress() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
