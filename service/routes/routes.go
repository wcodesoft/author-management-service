package routes

import (
	"context"
	authorGrpc "github.com/wcodesoft/author-management-service/grpc/go/author-management.proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"service/database"
)

// Routes Struct with implementation of the RPC endpoints.
type Routes struct {
	authorGrpc.UnimplementedAuthorManagementServer
	db database.DbConnector
}

// NewRoutes Creates a new server with the passed connector.
func NewRoutes(database database.DbConnector) *Routes {
	s := &Routes{
		db: database,
	}
	return s
}

// CreateAuthor Creates an author in the database from the received data.
func (server *Routes) CreateAuthor(_ context.Context, author *authorGrpc.Author) (*authorGrpc.Response, error) {
	_, err := server.db.AddAuthor(author.Uuid, author.Name, author.PicUrl)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}

// GetAuthor Gets an Author from the database and return to the user.
func (server *Routes) GetAuthor(_ context.Context, requestID *authorGrpc.RequestId) (*authorGrpc.Author, error) {
	author, err := server.db.GetAuthor(requestID.GetUuid())
	id := author.UUID.String()
	return &authorGrpc.Author{
		Uuid:   &id,
		Name:   author.Name,
		PicUrl: author.PicURL,
	}, err
}

// GetAuthors Retrieve all authors stored on the database.
func (server *Routes) GetAuthors(_ context.Context, _ *emptypb.Empty) (*authorGrpc.GetAuthorResponse, error) {
	allAuthors := server.db.GetAuthors()
	var array []*authorGrpc.Author

	for _, author := range allAuthors {
		authorUUID := author.UUID.String()
		var rpcAuthor = authorGrpc.Author{
			Uuid:   &authorUUID,
			Name:   author.Name,
			PicUrl: author.PicURL,
		}
		array = append(array, &rpcAuthor)
	}

	return &authorGrpc.GetAuthorResponse{
		Author: array,
	}, nil
}

// UpdateAuthor Update author entry with the new data.
func (server *Routes) UpdateAuthor(_ context.Context, author *authorGrpc.Author) (*authorGrpc.Response, error) {
	err := server.db.UpdateAuthor(*author.Uuid, author.Name, author.PicUrl)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}

// DeleteAuthor Delete an author entry from the database.
func (server *Routes) DeleteAuthor(_ context.Context, request *authorGrpc.RequestId) (*authorGrpc.Response, error) {
	err := server.db.DeleteAuthor(request.Uuid)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}
