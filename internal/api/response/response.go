package response

type ResponseHTTP struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
