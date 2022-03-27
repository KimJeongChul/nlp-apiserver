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

	"github.com/gorilla/mux"
)

/*
	Method: POST
	Request:
	 - Content-Type: application/json
	 - {
		 "modelType": "intent" or "ner"
	   }
	Response:
	 - Content-Type: application/json
	 - {
		  "trainID": "TRAIN-ID(UUID v5)"
          "resultCode": 200,
		  "resultMsg":  "success",
	   }
*/
// handleTrainCreate POST /train Generate trainID
func (as *ApiServer) handleTrainCreate(w http.ResponseWriter, req *http.Request) {
	session, cErr, resCode := as.APICallProcessing(&w, req, "train")
	if cErr != nil {
		logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
		as.GenResponseErrorMessage(w, cErr, resCode)
	}

	defer func() {
		if cErr != nil {
			logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
			as.GenResponseErrorMessage(w, cErr, resCode)
		}
	}()

	switch req.Method {
	case "POST":
		bodyLen := req.ContentLength
		body, errReadAll := ioutil.ReadAll(req.Body)
		if errReadAll != nil {
			cErr, resCode = errors.NewCError(errors.HTTP_BODY_READ_ERR, errReadAll.Error()), http.StatusInternalServerError
			return
		}
		logger.LogI(session.FuncName, session.TransactionId, "req bodylen: ", bodyLen, " body: ", string(body))

		var msgTrainCreateReq MsgTrainCreateReq
		errUnmarshal := json.Unmarshal(body, &msgTrainCreateReq)
		if errUnmarshal != nil {
			cErr, resCode = errors.NewCError(errors.JSON_UNMARSHAL_ERR, errUnmarshal.Error()), http.StatusBadRequest
			return
		}

		// Check valid for HTTP request body
		if msgTrainCreateReq.ModelType == nil {
			cErr, resCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) is null"), http.StatusBadRequest
			return 
		}
		modelType := *msgTrainCreateReq.ModelType
		if modelType == "" {
			cErr, resCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) value is empty"), http.StatusBadRequest
			return
		}
		if !as.CheckRequestModelType(modelType) {
			cErr, resCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) value need to set(intent, ner)"), http.StatusBadRequest
			return
		}

		trainID, errGenTrxId := as.getUniqueTrxId()
		if errGenTrxId != nil {
			cErr, resCode = errors.NewCError(errors.HTTP_PREPROCESSING_ERR, "Cannot Create TransactionId"), http.StatusInternalServerError
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

		trainIDDir := dataDir + "/" + trainID
		if _, errNotExist := os.Stat(trainIDDir); errorslib.Is(errNotExist, os.ErrNotExist) {
			logger.LogD(session.FuncName, session.TransactionId, "create directory:", trainIDDir)
			errMkdir := os.Mkdir(trainIDDir, 0700)
			if errMkdir != nil  {
				cErr, resCode = errors.NewCError(errors.OS_MKDIR_ERR, errMkdir.Error()), http.StatusInternalServerError
				return
			}
		}

		resMsg := MsgTrainCreateRes {
			TrainId: trainID,
			ResultCode: resCode,
			ResultMsg:  SUCCESS,
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

/*
	Method: POST
	Request: None
	Response:
	 - Content-Type: application/json
	 - {
          "resultCode": 200,
		  "resultMsg":  "success",
	   }
*/
// handleIntentModelDelete /train/{trainID}/delete 의도분류 모델 삭제
func (as *ApiServer) handleTrainDelete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	trainID := vars["trainID"]

	session, cErr, resCode := as.APICallProcessing(&w, req, "trainDelete")
	if cErr != nil {
		logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
		as.GenResponseErrorMessage(w, cErr, resCode)
	}

	defer func() {
		if cErr != nil {
			logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
			as.GenResponseErrorMessage(w, cErr, resCode)
		}
	}()

	switch req.Method {
	case "POST":
		dataDir := as.serverConfig.HomePath + "/data"
		if _, errNotExist := os.Stat(dataDir); errorslib.Is(errNotExist, os.ErrNotExist) {
			logger.LogD(session.FuncName, session.TransactionId, "create directory:", dataDir)
			errMkdir := os.Mkdir(dataDir, 0700)
			if errMkdir != nil  {
				cErr, resCode = errors.NewCError(errors.OS_MKDIR_ERR, errMkdir.Error()), http.StatusInternalServerError
				return
			}
		}

		trainIDDir := dataDir + "/" + trainID
		if _, errNotExist := os.Stat(trainIDDir); errorslib.Is(errNotExist, os.ErrNotExist) {
			msgErr := "Err Msg= " + trainIDDir + " directory not exist"
			logger.LogE(session.FuncName, session.TransactionId, msgErr)
			returnMsg := MsgErrorRes {
				ResultCode: 404,
				ResultMsg:  msgErr,
			}
			wByte, errMarshal := json.Marshal(returnMsg)
			if errMarshal != nil {
				cErr, resCode = errors.NewCError(errors.JSON_MARSHAL_ERR, errMarshal.Error()), http.StatusInternalServerError
				return
			}
			w.Write(wByte)
			cErr, resCode = nil, 200
			return
		}
		// Delete trainID directory
		errRemove := os.RemoveAll(trainIDDir)
		if errRemove != nil {
			msgErr := "Err Msg= remove " + trainIDDir + " error"
			returnMsg := MsgErrorRes{
				ResultCode: 500,
				ResultMsg:  msgErr,
			}
			wByte, err := json.Marshal(returnMsg)
			if err != nil {
				cErr, resCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
				return
			}
			w.Write(wByte)
			cErr, resCode = nil, 200
			return
		}

		resMsg := MsgSuccessRes {
			ResultCode: 200,
			ResultMsg:  SUCCESS,
		}
		wByte, err := json.Marshal(resMsg)
		if err != nil {
			wByte = resMsg.generateSuccMsg()
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