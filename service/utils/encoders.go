package utils

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventManager "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
)

// EncodeResponseToByte Encodes the proto Response into a byte array.
func EncodeResponseToByte(response *eventManager.Response) []byte {
	encoded, _ := proto.Marshal(response)
	return encoded
}

// EncodeAuthorToString Encodes the proto Author into a base64 serialized string.
func EncodeAuthorToString(author *authorManagementProto.Author) string {
	encoded, _ := proto.Marshal(author)
	encodedString := base64.StdEncoding.EncodeToString(encoded)
	return encodedString
}

// EncodeAuthorsListToString Encodes the proto AuthorList into a base64 serialized string.
func EncodeAuthorsListToString(list *authorManagementProto.AuthorList) string {
	encoded, _ := proto.Marshal(list)
	encodedString := base64.StdEncoding.EncodeToString(encoded)
	return encodedString
}
