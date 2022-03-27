package api

import (
	"encoding/json"
	"net/http"
	"nlp-apiserver/errors"
)

// GenResponseErrorMessage Response 에러 메시지 전송
func (as *ApiServer) GenResponseErrorMessage(w http.ResponseWriter, cErr *errors.CError, errCode int) string {
	w.Header().Set("Content-Type", "application/json")
	msgErrorRes := MsgErrorRes{
		ResultCode: errCode,
		ResultMsg:  cErr.Message,
	}
	wByte, err := json.Marshal(msgErrorRes)
	if err != nil {
		cErr, errCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
		message := "{\"resultCode\":500, \"result\":\""+cErr.Error()+"\"}"
		w.Write([]byte(message))
		return message
	}
	w.Write(wByte)
	return string(wByte)
}
