package api

import (
	"encoding/json"
	errorslib "errors"
	"io/ioutil"
	"net/http"
	"nlp-apiserver/errors"
	"nlp-apiserver/logger"
	"os"
	"time"
)

// handleTrainCreate POST /train Generate train_id
func (as *ApiServer) handleTrainCreate(w http.ResponseWriter, req *http.Request) {
	apiType := "train"
	var cErr *errors.CError
	var resCode int

	session, errApi := as.APICallProcessing(&w, req, apiType)
	if errApi != nil {
		logger.LogE(session.FuncName, session.TransactionId, "APICallProcess Err Msg=", errApi.Error())
		as.responseErrorMessage(w, errApi, http.StatusInternalServerError)
	}
	defer as.APICallPostprocessing(w, session, cErr, resCode)
	logger.LogI(session.FuncName, session.TransactionId, "/" + apiType +" API call")
	
	switch req.Method {
	case "POST":
		bodylen := req.ContentLength
		body, errIO := ioutil.ReadAll(req.Body)
		if errIO != nil {
			cErr, resCode = errors.NewCError(errors.HTTP_BODY_READ_ERR, errIO.Error()), http.StatusInternalServerError
			return
		}
		logger.LogI(session.FuncName, session.TransactionId, "req bodylen: ", bodylen, " body: ", string(body))

		var msgTrainReq MsgTrainReq
		if errUnmarshal := json.Unmarshal(body, &msgTrainReq); errUnmarshal != nil {
			cErr, resCode = errors.NewCError(errors.JSON_UNMARSHAL_ERR, errUnmarshal.Error()), http.StatusBadRequest
			return
		}

		// Check parameter validation
		if msgTrainReq.ModelType == nil {
			cErr, resCode = errors.NewCError(errors.BAD_REQUEST_PARAMETER_ERR, "model_type parameter is not exist"), http.StatusBadRequest
			return
		}
		modelType := *msgTrainReq.ModelType
		if cErr, resCode = as.checkValidModelType(modelType); cErr != nil {
			return
		}

		// Create Train ID
		trainId, errTrxId := as.getUniqueTrxId()
		if errTrxId != nil {
			cErr, resCode = errors.NewCError(errors.GENERATE_TRANSACION_ID_ERR, errTrxId.Error()), http.StatusInternalServerError
			return
		}

		dataDir := as.serverConfig.HomePath + "/data"
		if _, errNotExist := os.Stat(dataDir); errorslib.Is(errNotExist, os.ErrNotExist) {
			logger.LogD(session.FuncName, session.TransactionId, "create directory:", dataDir)
			errMkdir := os.Mkdir(dataDir, 0700)
			if errMkdir != nil  {
				cErr, resCode = errors.NewCError(errors.OS_MKDIR_ERR, errMkdir.Error()), http.StatusInternalServerError
				return
			}
		}

		trainIdDir := dataDir + "/" + trainId
		if _, errNotExist := os.Stat(trainIdDir); errorslib.Is(errNotExist, os.ErrNotExist) {
			logger.LogD(session.FuncName, session.TransactionId, "create directory:", trainIdDir)
			errMkdir := os.Mkdir(trainIdDir, 0700)
			if errMkdir != nil  {
				cErr, resCode = errors.NewCError(errors.OS_MKDIR_ERR, errMkdir.Error()), http.StatusInternalServerError
				return
			}
		}

		resMsg := MsgTrainRes {
			TrainId: trainId,
			ResultCode: resCode,
			Result: SUCCESS,
		}
		wByte, err := json.Marshal(resMsg)
		if err != nil {
			wByte = resMsg.generateTrainSuccMsg()
		}
		w.Write(wByte)

		// Check response time
		responseTime := time.Now().Sub(session.StartResponseTime)
		logger.LogI(session.FuncName, session.TransactionId, "/" + session.ApiType + " API call finish. responseTime=", responseTime)
	
	default:
		cErr, resCode = errors.NewCError(errors.HTTP_INVALID_METHOD_ERR, "Invalid Method"), http.StatusMethodNotAllowed
		return
	}
}