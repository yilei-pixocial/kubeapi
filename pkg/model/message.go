package model

type Message struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Ok(resp interface{}) Message {
	result := Message{Code: MESSAGE_OK, Message: "success", Result: resp}
	return result
}

func Error(msg string) Message {
	return Message{Code: MESSAGE_ERROR, Message: msg, Result: nil}
}
