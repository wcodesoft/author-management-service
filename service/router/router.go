package router

import (
	"errors"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"service/database"
	"service/utils"
)

type RouteManager struct {
	connector database.DbConnector
}

func NewRouteManager(connector database.DbConnector) *RouteManager {
	return &RouteManager{
		connector: connector,
	}
}

func (rm *RouteManager) RouteEvent(event *eventProto.Event) ([]string, error) {
	switch event.Action {
	case eventProto.Action_CREATE:
		return rm.createAuthor(event)
	case eventProto.Action_UPDATE:
		return rm.updateAuthor(event)
	case eventProto.Action_READ:
		return rm.readAuthor(event)
	case eventProto.Action_DELETE:
		return rm.deleteAuthor(event)
	}
	return nil, errors.New("action not supported")
}

func (rm *RouteManager) createAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	uuid, err := rm.connector.AddAuthor(database.AuthorFromGrpc(author))
	parsedUuid := uuid.String()
	return []string{parsedUuid}, err
}

func (rm *RouteManager) updateAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	err := rm.connector.UpdateAuthor(database.AuthorFromGrpc(author))
	return nil, err
}

func (rm *RouteManager) readAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.AllEntries {
		return rm.readAllAuthors()
	} else {
		return rm.readAuthorById(query.GetUuid())
	}
}

func (rm *RouteManager) readAuthorById(uuid string) ([]string, error) {
	author, err := rm.connector.GetAuthor(uuid)
	parsedAuthor := database.AuthorToGrpc(*author)
	return []string{utils.EncodeAuthorToString(parsedAuthor)}, err
}

func (rm *RouteManager) readAllAuthors() ([]string, error) {
	authors := rm.connector.GetAuthors()
	parsedAuthors := database.AuthorListToGrpcList(authors)
	return []string{utils.EncodeAuthorsListToString(&parsedAuthors)}, nil
}

func (rm *RouteManager) deleteAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.Uuid == nil {
		return nil, errors.New("uuid not set on the request")
	}

	err := rm.connector.DeleteAuthor(query.GetUuid())
	return nil, err
}
