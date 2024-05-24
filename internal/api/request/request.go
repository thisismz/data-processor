package request

type DataRequest struct {
	DataID  string `json:"data_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
	Payload []byte `json:"payload" binding:"required"`
}
