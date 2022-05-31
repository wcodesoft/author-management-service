package routes

import (
	"context"
	authorGrpc "github.com/wcodesoft/author-management-service/grpc/go/author-management.proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"service/database"
)

type routeServer struct {
	authorGrpc.UnimplementedAuthorManagementServer
	db database.Database
}

func NewRouteServer(database database.Database) *routeServer {
	s := &routeServer{
		db: database,
	}
	return s
}

func (server *routeServer) CreateAuthor(_ context.Context, author *authorGrpc.Author) (*authorGrpc.Response, error) {
	_, err := server.db.AddAuthor(author.Uuid, author.Name, author.PicUrl)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}

func (server *routeServer) GetAuthor(_ context.Context, requestId *authorGrpc.RequestId) (*authorGrpc.Author, error) {
	author, err := server.db.GetAuthor(requestId.GetUuid())
	id := author.UUID.String()
	return &authorGrpc.Author{
		Uuid:   &id,
		Name:   author.Name,
		PicUrl: author.PicURL,
	}, err
}

func (server *routeServer) GetAuthors(_ context.Context, _ *emptypb.Empty) (*authorGrpc.GetAuthorResponse, error) {
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

func (server *routeServer) UpdateAuthor(_ context.Context, author *authorGrpc.Author) (*authorGrpc.Response, error) {
	err := server.db.UpdateAuthor(*author.Uuid, author.Name, author.PicUrl)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}

func (server *routeServer) DeleteAuthor(_ context.Context, request *authorGrpc.RequestId) (*authorGrpc.Response, error) {
	err := server.db.DeleteAuthor(request.Uuid)
	return &authorGrpc.Response{
		Success: err == nil,
	}, err
}
