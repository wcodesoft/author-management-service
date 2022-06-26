package utils

import (
	"encoding/base64"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"google.golang.org/protobuf/proto"
)

func DecodeEvent(body []byte) *eventProto.Event {
	event := &eventProto.Event{}
	proto.Unmarshal(body, event)
	return event
}

func DecodeQuery(message string) *eventProto.Query {
	decoded, _ := base64.StdEncoding.DecodeString(message)
	query := &eventProto.Query{}
	proto.Unmarshal(decoded, query)
	return query
}

func DecodeAuthor(message string) *authorManagementProto.Author {
	decoded, _ := base64.StdEncoding.DecodeString(message)
	author := &authorManagementProto.Author{}
	proto.Unmarshal(decoded, author)
	return author
}
