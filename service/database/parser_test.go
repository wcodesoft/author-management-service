package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	"testing"
)

func TestAuthorFromGrpc(t *testing.T) {
	expectedUUID := uuid.New().String()
	expectedName := "Test"
	expectedPic := "TestPic"
	authorGrpc := &authorManagementProto.Author{
		Uuid:   &expectedUUID,
		Name:   expectedName,
		PicUrl: &expectedPic,
	}
	parsedAuthor := AuthorFromGrpc(authorGrpc)
	picURL := *parsedAuthor.PicURL
	assert.Equal(t, expectedUUID, parsedAuthor.ID.String())
	assert.Equal(t, expectedName, parsedAuthor.Name)
	assert.Equal(t, expectedPic, picURL)
}

func TestAuthorListToGrpcList(t *testing.T) {
	var authorList []Author
	expectedLen := 3
	for i := 0; i < expectedLen; i++ {
		newUUID := uuid.New()
		authorList = append(authorList, Author{
			ID:     &newUUID,
			Name:   "Test",
			PicURL: nil,
		})
	}
	parsedAuthors := AuthorListToGrpcList(authorList)
	assert.Len(t, parsedAuthors.Authors, expectedLen)
}

func TestAuthorToGrpc(t *testing.T) {
	newUUID := uuid.New()
	author := Author{
		ID:     &newUUID,
		Name:   "John Doe",
		PicURL: nil,
	}
	grpcAuthor := AuthorToGrpc(author)
	authorIDString := author.ID.String()
	assert.Equal(t, author.Name, grpcAuthor.Name)
	assert.Equal(t, &authorIDString, grpcAuthor.Uuid)
	assert.Equal(t, author.PicURL, grpcAuthor.PicUrl)
}
