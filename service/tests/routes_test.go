package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	authorManagement "github.com/wcodesoft/author-management-service/grpc/go/author-management.proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"service/database"
	"service/routes"
	"testing"
)

func localDatabase() database.Database {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}
	db.AutoMigrate(&database.Author{})
	return database.Database{
		Database: db,
	}
}

func getNewUUID() string {
	newUUID, _ := uuid.NewUUID()
	uuidString := newUUID.String()
	return uuidString
}

func TestCreateAuthorRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Name:   name,
		PicUrl: &picUrl,
	}

	resp, err := server.CreateAuthor(ctx, &author)

	assert.NoError(t, err, "Error when creating new author.")
	assert.Equal(t, resp.Success, true)
}

func TestCreateAuthorWithUUIDRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	newUUID := getNewUUID()
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Uuid:   &newUUID,
		Name:   name,
		PicUrl: &picUrl,
	}

	resp, err := server.CreateAuthor(ctx, &author)

	assert.NoError(t, err, "Error when creating new author.")
	assert.Equal(t, resp.Success, true)

	authorResponse, err := server.GetAuthor(ctx, &authorManagement.RequestId{
		Uuid: newUUID,
	})

	assert.NoError(t, err)
	assert.Equal(t, newUUID, authorResponse.GetUuid())
	assert.Equal(t, name, authorResponse.GetName())
	assert.Equal(t, picUrl, authorResponse.GetPicUrl())
}

func TestGetUserRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	newUUID := getNewUUID()
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Uuid:   &newUUID,
		Name:   name,
		PicUrl: &picUrl,
	}

	resp, _ := server.CreateAuthor(ctx, &author)
	assert.Equal(t, resp.Success, true)

	authorResponse, err := server.GetAuthor(ctx, &authorManagement.RequestId{
		Uuid: newUUID,
	})

	assert.NoError(t, err)
	assert.Equal(t, newUUID, authorResponse.GetUuid())
	assert.Equal(t, name, authorResponse.GetName())
	assert.Equal(t, picUrl, authorResponse.GetPicUrl())
}

func TestDeleteUserRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	newUUID := getNewUUID()
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Uuid:   &newUUID,
		Name:   name,
		PicUrl: &picUrl,
	}

	resp, _ := server.CreateAuthor(ctx, &author)
	assert.Equal(t, resp.Success, true)

	deleteResp, err := server.DeleteAuthor(ctx, &authorManagement.RequestId{
		Uuid: newUUID,
	})

	assert.NoError(t, err)
	assert.True(t, deleteResp.Success)
}

func TestUpdateUserRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	newUUID := getNewUUID()
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Uuid:   &newUUID,
		Name:   name,
		PicUrl: &picUrl,
	}

	resp, _ := server.CreateAuthor(ctx, &author)
	assert.Equal(t, resp.Success, true)

	newName := "New Author Name"

	updateResp, err := server.UpdateAuthor(ctx, &authorManagement.Author{
		Uuid:   &newUUID,
		Name:   newName,
		PicUrl: &picUrl,
	})

	assert.NoError(t, err)
	assert.True(t, updateResp.Success)

	authorResponse, err := server.GetAuthor(ctx, &authorManagement.RequestId{
		Uuid: newUUID,
	})

	assert.NoError(t, err)
	assert.Equal(t, newUUID, authorResponse.GetUuid())
	assert.Equal(t, newName, authorResponse.GetName())
	assert.Equal(t, picUrl, authorResponse.GetPicUrl())
}

func TestGetAuthorsRPC(t *testing.T) {
	db := localDatabase()
	ctx := context.Background()
	defer db.CloseDatabase()

	server := routes.NewRouteServer(db)
	newUUID := getNewUUID()
	name := "Author Name"
	picUrl := "picture"

	author := authorManagement.Author{
		Uuid:   &newUUID,
		Name:   name,
		PicUrl: &picUrl,
	}

	name2 := "Second name"

	author2 := authorManagement.Author{
		Name: name2,
	}

	resp, _ := server.CreateAuthor(ctx, &author)
	assert.Equal(t, resp.Success, true)
	resp2, _ := server.CreateAuthor(ctx, &author2)
	assert.Equal(t, resp2.Success, true)

	authors, err := server.GetAuthors(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, len(authors.Author), 2)
}
