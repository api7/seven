package apisix

type ErrorResponse struct {
	Cause     string `json:"cause"`
	Index     int64  `json:"index"`
	ErrorCode int64  `json:"errorCode"`
	Message   string `json:"message"`
}
