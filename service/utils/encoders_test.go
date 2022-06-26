package utils

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventManagerProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"testing"
)

func TestEncodeResponseToByte(t *testing.T) {
	expectedResponse := &eventManagerProto.Response{
		Success: true,
		Error:   nil,
	}
	byteResponse := EncodeResponseToByte(expectedResponse)
	resultResponse := &eventManagerProto.Response{}
	err := proto.Unmarshal(byteResponse, resultResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Success, resultResponse.Success)
}

func TestEncodeAuthorToString(t *testing.T) {
	newUuid := uuid.NewString()
	author := &authorManagementProto.Author{
		Uuid: &newUuid,
		Name: "Walter José",
	}
	encoded, _ := proto.Marshal(author)
	expectedBase64 := base64.StdEncoding.EncodeToString(encoded)
	resultString := EncodeAuthorToString(author)
	assert.Equal(t, expectedBase64, resultString)
}

func TestEncodeAuthorsListToString(t *testing.T) {
	var authors []*authorManagementProto.Author
	expectedLen := 3
	for i := 0; i < expectedLen; i++ {
		newUuid := uuid.NewString()
		authors = append(authors, &authorManagementProto.Author{
			Uuid: &newUuid,
			Name: "Test",
		})
	}
	authorsList := &authorManagementProto.AuthorList{
		Authors: authors,
	}
	encoded, _ := proto.Marshal(authorsList)
	expectedBase64 := base64.StdEncoding.EncodeToString(encoded)
	resultString := EncodeAuthorsListToString(authorsList)
	assert.Equal(t, expectedBase64, resultString)
}
