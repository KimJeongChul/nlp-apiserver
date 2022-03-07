package api

import (
	"encoding/json"
	"net/http"
	"nlp-apiserver/errors"
)

const (
	SUCCESS = "success"
)

type MsgCommonRes struct {
	ResultCode int    `json:"result_code"`
	Result     string `json:"result"`
}

func (as *ApiServer) responseErrorMessage(w http.ResponseWriter, cErr *errors.CError, errCode int) {
	msgCommonRes := MsgCommonRes{
		ResultCode: errCode,
		Result:     cErr.Message,
	}
	wByte, err := json.Marshal(msgCommonRes)
	if err != nil {
		cErr, errCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
		message := "{\"resultCode\":500, \"result\":\"" + cErr.Error() + "\"}"
		w.Write([]byte(message))
	}
	w.Write(wByte)
}

// MsgTrainReq /train Request
type MsgTrainReq struct {
	ModelType *string `json:"model_type"`
}

type MsgTrainRes struct {
	TrainId    string `json:"train_id"`
	ResultCode int    `json:"result_code"`
	Result     string `json:"result"`
}

func (mtr MsgTrainRes) generateTrainSuccMsg() []byte {
	return []byte("{\"resultCode\":200, \"result\":\"success\", \"train_id\":\"" + mtr.TrainId+"\"}")
}