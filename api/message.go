package api

type MsgErrorRes struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

type MsgTrainCreateReq struct {
	ModelType *string `json:"modelType"`
}
