package utils

import (
	"encoding/base64"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"google.golang.org/protobuf/proto"
)

// DecodeEvent Receives an array of bytes and transform to proto Event.
func DecodeEvent(body []byte) *eventProto.Event {
	event := &eventProto.Event{}
	proto.Unmarshal(body, event)
	return event
}

// DecodeQuery Receives a base64 serialized string and parse it to a proto Query.
func DecodeQuery(message string) *eventProto.Query {
	decoded, _ := base64.StdEncoding.DecodeString(message)
	query := &eventProto.Query{}
	proto.Unmarshal(decoded, query)
	return query
}

// DecodeAuthor Receives a base64 serialized string and parse it to a proto Author.
func DecodeAuthor(message string) *authorManagementProto.Author {
	decoded, _ := base64.StdEncoding.DecodeString(message)
	author := &authorManagementProto.Author{}
	proto.Unmarshal(decoded, author)
	return author
}
