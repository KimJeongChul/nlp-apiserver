package api

const (
	SUCCESS = "success"
)

type MsgErrorRes struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

type MsgSuccessRes struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

func (msr MsgSuccessRes) generateSuccMsg() []byte {
	return []byte("{\"resultCode\":200, \"result\":\"success\"}")
}

type MsgTrainCreateReq struct {
	ModelType *string `json:"modelType"`
}

type MsgTrainCreateRes struct {
	TrainId    string `json:"trainID"`
	ResultCode int    `json:"resultCode"`
	ResultMsg     string `json:"resultMsg"`
}

func (mtcr MsgTrainCreateRes) generateTrainSuccMsg() []byte {
	return []byte("{\"resultCode\":200, \"result\":\"success\", \"train_id\":\"" + mtcr.TrainId+"\"}")
}