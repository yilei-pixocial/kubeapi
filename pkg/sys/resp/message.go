package resp

type Message struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BadResponse(msg string) Message {
	return Message{Code: MESSAGE_BAD, Message: msg}
}

func Ok() Message {
	result := Message{Code: MESSAGE_OK, Message: "ok"}
	return result
}

func OkWithData(data interface{}) Message {
	return Message{Code: MESSAGE_OK, Message: "success", Data: data}
}

func ErrorWithMsg(msg string) Message {
	return Message{Code: MESSAGE_ERROR, Message: msg}
}
