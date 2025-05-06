package common

type MessageError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
