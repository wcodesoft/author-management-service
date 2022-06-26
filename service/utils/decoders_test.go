package utils

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"testing"
)

func TestDecodeEvent(t *testing.T) {
	expectedEvent := &eventProto.Event{
		Action:  eventProto.Action_READ,
		Message: "EventProto",
	}
	encoded, _ := proto.Marshal(expectedEvent)
	decoded := DecodeEvent(encoded)
	assert.Equal(t, expectedEvent.Action, decoded.Action)
	assert.Equal(t, expectedEvent.Message, decoded.Message)
}

func TestDecodeQuery(t *testing.T) {
	uuidString := uuid.NewString()
	expectedQuery := &eventProto.Query{
		AllEntries: false,
		Uuid:       &uuidString,
	}
	encoded, _ := proto.Marshal(expectedQuery)
	queryString := base64.StdEncoding.EncodeToString(encoded)
	decodedQuery := DecodeQuery(queryString)
	assert.Equal(t, expectedQuery.Uuid, decodedQuery.Uuid)
	assert.Equal(t, expectedQuery.AllEntries, decodedQuery.AllEntries)
}

func TestDecodeAuthor(t *testing.T) {
	uuidString := uuid.NewString()
	expectedAuthor := &authorManagementProto.Author{
		Uuid: &uuidString,
		Name: "John Doe",
	}
	encoded, _ := proto.Marshal(expectedAuthor)
	authorString := base64.StdEncoding.EncodeToString(encoded)
	decodedAuthor := DecodeAuthor(authorString)
	assert.Equal(t, expectedAuthor.Name, decodedAuthor.Name)
	assert.Equal(t, expectedAuthor.Uuid, decodedAuthor.Uuid)
	assert.Equal(t, expectedAuthor.PicUrl, decodedAuthor.PicUrl)
}
