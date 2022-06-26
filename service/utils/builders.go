package utils

import eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"

// BuildResponse Receives an array of string and an error and parse them into a Response.
func BuildResponse(message []string, err error) *eventProto.Response {
	if err != nil {
		errorString := err.Error()
		return &eventProto.Response{
			Success: false,
			Error:   &errorString,
		}
	}
	return &eventProto.Response{
		Success: true,
		Result:  message,
	}
}
