package api

import (
	"net/http"
	"nlp-apiserver/errors"
	"nlp-apiserver/logger"
	"nlp-apiserver/utils"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type SessionObject struct {
	FuncName          string
	TransactionId     string
	StartResponseTime time.Time
}

func (as *ApiServer) APICallProcessing(w *http.ResponseWriter, req *http.Request) (sessionObj *SessionObject, cErr *errors.CError) {
	var err error
	sessionObj = &SessionObject{}

	//Enable CORS
	utils.EnableCors(w)

	// Check start response time
	sessionObj.StartResponseTime = time.Now()

	sessionObj.TransactionId, err = as.getUniqueTrxId()
	if err != nil {
		cErr = errors.NewCError(errors.HTTP_PREPROCESSING_ERR, "Cannot Create TransactionId")
		logger.LogE(sessionObj.FuncName, sessionObj.TransactionId, "ERROR:Msg=", cErr.Error())
		return
	}

	//Set Function Name
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		sessionObj.FuncName = "UNKNOWN"
	} else {
		funcNameArr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		sessionObj.FuncName = funcNameArr[len(funcNameArr)-1]
	}

	return
}

// GetUniqueTrxId Generate UUID
func (as *ApiServer) getUniqueTrxId() (string, error) {
	uuidCnt = (uuidCnt + 1) % 10000
	uuidString := utils.GetMillisTimeFormat(time.Now()) + ":" + localMacAddr + ":" + strconv.Itoa(uuidCnt)
	transactionId := uuid.NewV5(uuid.NamespaceDNS, uuidString)
	return transactionId.String(), nil
}