package api

import (
	"encoding/json"
	"net/http"
	"nlp-apiserver/errors"
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
