package tests

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"service/database"
	"testing"
)

func TestAddAuthorNotPassingUUID(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	_, err := db.AddAuthor(nil, "John Doe", nil)
	assert.NoError(t, err, "Fail when adding user.")
}

func TestAddAuthorPassingUUID(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	newUUID, _ := uuid.NewUUID()
	newUUIDString := newUUID.String()
	id, err := db.AddAuthor(&newUUIDString, "John Doe", nil)
	assert.NoError(t, err, "Fail when adding user.")
	assert.Equal(t, newUUIDString, id.String())
}

func TestAddAuthorPassingInvalidUUID(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	newUUID := "NotValidUUID"
	id, err := db.AddAuthor(&newUUID, "John Doe", nil)
	assert.Nil(t, id, "Adding with invalid UUID")
	assert.Errorf(t, err, "Adding with invalid UUID")
}

func TestGetAuthor(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	picUrl := "johndoe"
	ans, err := db.AddAuthor(nil, "John Doe", &picUrl)
	assert.NoError(t, err, "Fail when adding user.")
	author, errGet := db.GetAuthor(ans.String())
	assert.NoError(t, errGet, "Fail when retrieving author")
	assert.Equal(t, "John Doe", author.Name)
	assert.Equal(t, "johndoe", *author.PicURL)
}

func TestGetAllAuthors(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	db.AddAuthor(nil, "Author1", nil)
	authors := db.GetAuthors()
	assert.Len(t, authors, 1, "Wrong number of authors, expected 1 got %d", len(authors))
	db.AddAuthor(nil, "Author2", nil)
	authors = db.GetAuthors()
	assert.Len(t, authors, 2, "Wrong number of authors, expected 1 got %d", len(authors))
}

func TestUpdateAuthor(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	authorId, err := db.AddAuthor(nil, "Author1", nil)
	assert.NoError(t, err, "Fail to add an author.")
	newPicUrl := "newPicUrl"
	err = db.UpdateAuthor(authorId.String(), "Author1", &newPicUrl)
	assert.NoError(t, err, "Fail to update author data.")
	var author, errGet = db.GetAuthor(authorId.String())
	assert.NoError(t, errGet, "Fail to get author")
	assert.Equal(t, author.Name, "Author1")
	assert.Equal(t, *author.PicURL, "newPicUrl")
}

func TestUpdateAuthorNonExistentAuthor(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	err := db.UpdateAuthor("NotExistentUUID", "", nil)
	assert.Error(t, err, "Able to update author.")
}

func TestDeleteAuthor(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	authorId, err := db.AddAuthor(nil, "Author1", nil)
	assert.NoError(t, err, "Fail to add an author.")
	err = db.DeleteAuthor(authorId.String())
	assert.NoError(t, err, "Fail to delete author data.")
	var author, errGet = db.GetAuthor(authorId.String())
	assert.Error(t, errGet, "Not able to get author data because was deleted.")
	assert.Nil(t, author, "Author was not deleted but retrieved.")
}

func TestDeleteAuthorNonExistentAuthor(t *testing.T) {
	db := database.NewDatabase()
	defer db.CloseDatabase()
	err := db.DeleteAuthor("NonExistentUUID")
	assert.Error(t, err, "Able to delete entry.")
}
