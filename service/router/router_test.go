package router

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	authorManagementProto "github.com/wcodesoft/author-management-service/protos/go/author-management.proto"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"gorm.io/driver/sqlite"
	"service/database"
	"service/utils"
	"testing"
)

func TestRouteManager_CreateEvent(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])
}

func TestRouteManager_UpdateEvent(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])

	newAuthor := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	newAuthorString := utils.EncodeAuthorToString(&newAuthor)
	updateEvent := eventProto.Event{
		Action:  eventProto.Action_UPDATE,
		Message: newAuthorString,
	}

	result, err = router.RouteEvent(&updateEvent)
	assert.Nil(t, result)
	assert.NoError(t, err)
}

func TestRouteManager_ReadEvent(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])

	query := eventProto.Query{
		Uuid:       &newUuid,
		AllEntries: false,
	}
	byteQuery, _ := proto.Marshal(&query)
	queryString := base64.StdEncoding.EncodeToString(byteQuery)
	readEvent := eventProto.Event{
		Action:  eventProto.Action_READ,
		Message: queryString,
	}

	result, err = router.RouteEvent(&readEvent)
	receivedAuthor := utils.DecodeAuthor(result[0])
	assert.NoError(t, err)
	assert.Equal(t, author.Name, receivedAuthor.Name)
	assert.Equal(t, author.Uuid, receivedAuthor.Uuid)
	assert.Equal(t, author.PicUrl, receivedAuthor.PicUrl)
}

func TestRouteManager_ReadAllEvent(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])

	query := eventProto.Query{
		AllEntries: true,
	}
	byteQuery, _ := proto.Marshal(&query)
	queryString := base64.StdEncoding.EncodeToString(byteQuery)
	readEvent := eventProto.Event{
		Action:  eventProto.Action_READ,
		Message: queryString,
	}

	result, err = router.RouteEvent(&readEvent)
	decoded, _ := base64.StdEncoding.DecodeString(result[0])
	authorList := &authorManagementProto.AuthorList{}
	proto.Unmarshal(decoded, authorList)
	assert.NoError(t, err)
	authorsLen := len(authorList.Authors)
	assert.Equal(t, 1, authorsLen)
}

func TestRouteManager_DeleteEvent(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])

	query := eventProto.Query{
		Uuid:       &newUuid,
		AllEntries: false,
	}
	byteQuery, _ := proto.Marshal(&query)
	queryString := base64.StdEncoding.EncodeToString(byteQuery)
	event = eventProto.Event{
		Action:  eventProto.Action_DELETE,
		Message: queryString,
	}
	result, err = router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestRouteManager_DeleteEventWithoutUUID(t *testing.T) {
	sqliteDialector := sqlite.Open("file::memory:?cache=shared")
	db := database.NewConnection(sqliteDialector)
	defer db.CloseDatabase()
	newUuid := uuid.NewString()
	author := authorManagementProto.Author{
		Uuid:   &newUuid,
		Name:   "John Doe",
		PicUrl: nil,
	}
	authorString := utils.EncodeAuthorToString(&author)
	event := eventProto.Event{
		Action:  eventProto.Action_CREATE,
		Message: authorString,
	}

	router := NewRouteManager(db)
	result, err := router.RouteEvent(&event)
	assert.NoError(t, err)
	assert.Equal(t, newUuid, result[0])

	query := eventProto.Query{
		AllEntries: false,
	}
	byteQuery, _ := proto.Marshal(&query)
	queryString := base64.StdEncoding.EncodeToString(byteQuery)
	event = eventProto.Event{
		Action:  eventProto.Action_DELETE,
		Message: queryString,
	}
	result, err = router.RouteEvent(&event)
	assert.Error(t, err)
	assert.Nil(t, result)
}
