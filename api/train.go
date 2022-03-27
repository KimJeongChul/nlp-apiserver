package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"nlp-apiserver/errors"
	"nlp-apiserver/logger"
)

// handleTrainCreate Create Train ID
func (as *ApiServer) handleTrainCreate(w http.ResponseWriter, req *http.Request) {
	session, cErr, errCode := as.APICallProcessing(&w, req)
	if cErr != nil {
		logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
		as.GenResponseErrorMessage(w, cErr, errCode)
	}

	defer func() {
		if cErr != nil {
			logger.LogE(session.FuncName,  session.TransactionId, "Err Msg=", cErr.Error())
			as.GenResponseErrorMessage(w, cErr, errCode)
		}
	}()

	switch req.Method {
	case "POST":
		bodyLen := req.ContentLength
		body, errReadAll := ioutil.ReadAll(req.Body)
		if errReadAll != nil {
			cErr, errCode = errors.NewCError(errors.HTTP_BODY_READ_ERR, errReadAll.Error()), http.StatusInternalServerError
			return
		}
		logger.LogI(session.FuncName, session.TransactionId, "req bodylen: ", bodyLen, " body: ", string(body))

		var msgTrainCreateReq MsgTrainCreateReq
		errUnmarshal := json.Unmarshal(body, &msgTrainCreateReq)
		if errUnmarshal != nil {
			cErr, errCode = errors.NewCError(errors.JSON_UNMARSHAL_ERR, errUnmarshal.Error()), http.StatusBadRequest
			return
		}

		// Check valid for HTTP request body
		if msgTrainCreateReq.ModelType == nil {
			cErr, errCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) is null"), http.StatusBadRequest
			return 
		}
		modelType := *msgTrainCreateReq.ModelType
		if modelType == "" {
			cErr, errCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) value is empty"), http.StatusBadRequest
			return
		}
		if !as.CheckRequestModelType(modelType) {
			cErr, errCode = errors.NewCError(errors.INVALID_REQ_BODY_ERR, "Request body(modelType) value need to set(intent, ner)"), http.StatusBadRequest
			return
		}

		trainId, errGenTrxId = as.getUniqueTrxId()
		if errGenTrxId != nil {
			cErr, errCode = errors.NewCError(errors.HTTP_PREPROCESSING_ERR, "Cannot Create TransactionId"), http.StatusInternalServerError
			return
		}
		
		rootDataDir := ras.serverConfig.HomePath + "/data"


	}
}