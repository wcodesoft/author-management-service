package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"testing"
)

func TestAddAuthorNotPassingUUID(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author := Author{
		ID:     nil,
		Name:   "John Doe",
		PicURL: nil,
	}
	_, err := db.AddAuthor(author)
	assert.NoError(t, err, "Fail when adding user.")
}

func TestAddAuthorPassingUUID(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.New()
	author := Author{
		ID:     &newUuid,
		Name:   "John Doe",
		PicURL: nil,
	}
	newUUIDString := newUuid.String()
	id, err := db.AddAuthor(author)
	assert.NoError(t, err, "Fail when adding user.")
	assert.Equal(t, newUUIDString, id.String())
}

func TestAddAlreadyExistingAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author := Author{
		ID:     nil,
		Name:   "John Doe",
		PicURL: nil,
	}
	id, err := db.AddAuthor(author)
	assert.NoError(t, err, "Fail when adding user.")
	author2 := Author{
		ID:     id,
		Name:   "John Doe 2",
		PicURL: nil,
	}
	_, err2 := db.AddAuthor(author2)
	authors := db.GetAuthors()
	assert.Error(t, err2, "Author added with same Uuid %s", authors)

}

func TestGetAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	picUrl := "johndoe"
	newAuthor := Author{
		ID:     nil,
		Name:   "John Doe",
		PicURL: &picUrl,
	}
	ans, err := db.AddAuthor(newAuthor)
	assert.NoError(t, err, "Fail when adding user.")
	uuidString := ans.String()
	author, errGet := db.GetAuthor(uuidString)
	assert.NoError(t, errGet, "Fail when retrieving author")
	assert.Equal(t, "John Doe", author.Name)
	assert.Equal(t, "johndoe", *author.PicURL)
}

func TestGetAllAuthors(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author1 := Author{
		ID:     nil,
		Name:   "Author1",
		PicURL: nil,
	}
	db.AddAuthor(author1)
	authors := db.GetAuthors()
	assert.Len(t, authors, 1, "Wrong number of authors, expected 1 got %d", len(authors))
	author2 := Author{
		ID:     nil,
		Name:   "Author1",
		PicURL: nil,
	}
	db.AddAuthor(author2)
	authors = db.GetAuthors()
	assert.Len(t, authors, 2, "Wrong number of authors, expected 1 got %d", len(authors))
}

func TestUpdateAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author1 := Author{
		ID:     nil,
		Name:   "Author1",
		PicURL: nil,
	}
	authorId, err := db.AddAuthor(author1)
	assert.NoError(t, err, "Fail to add an author.")
	newPicUrl := "newPicUrl"
	newAuthor1 := Author{
		ID:     authorId,
		Name:   "Author1",
		PicURL: &newPicUrl,
	}
	err = db.UpdateAuthor(newAuthor1)
	assert.NoError(t, err, "Fail to update author data.")
	var author, errGet = db.GetAuthor(authorId.String())
	assert.NoError(t, errGet, "Fail to get author")
	assert.Equal(t, author.Name, "Author1")
	assert.Equal(t, *author.PicURL, "newPicUrl")
}

func TestUpdateAuthorWithNilUUID(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author1 := Author{
		ID:     nil,
		Name:   "Author1",
		PicURL: nil,
	}
	err := db.UpdateAuthor(author1)
	assert.Error(t, err, "Able to update author.")
}

func TestUpdateAuthorNonExistentAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUUID := uuid.New()
	author1 := Author{
		ID:     &newUUID,
		Name:   "Author1",
		PicURL: nil,
	}
	err := db.UpdateAuthor(author1)
	assert.Error(t, err, "Able to update author.")
}

func TestDeleteAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	author1 := Author{
		ID:     nil,
		Name:   "Author1",
		PicURL: nil,
	}
	authorId, err := db.AddAuthor(author1)
	assert.NoError(t, err, "Fail to add an author.")
	err = db.DeleteAuthor(authorId.String())
	assert.NoError(t, err, "Fail to delete author data.")
	var author, errGet = db.GetAuthor(authorId.String())
	assert.Error(t, errGet, "Not able to get author data because was deleted.")
	assert.Nil(t, author, "Author was not deleted but retrieved.")
}

func TestDeleteAuthorNonExistentAuthor(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	err := db.DeleteAuthor("NonExistentUUID")
	assert.Error(t, err, "Able to delete entry.")
}

func TestUUIDCorrectingParsing(t *testing.T) {
	newUuid := uuid.New()
	parsed := uuidParseOrCreate(newUuid.String())
	assert.Equal(t, newUuid.String(), parsed.String())
}

func TestUUIDCreatingNewUuidFromInvalidEntry(t *testing.T) {
	parsed := uuidParseOrCreate("Invalid")
	assert.NotEqual(t, parsed.String(), "Invalid")
}
