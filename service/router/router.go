package router

import (
	"errors"
	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
	"service/database"
	"service/utils"
)

// RouteManager manages the routing of received events.
type RouteManager struct {
	connector database.DbConnector
}

// NewRouteManager creates a new RouteManager instance.
func NewRouteManager(connector database.DbConnector) *RouteManager {
	return &RouteManager{
		connector: connector,
	}
}

// RouteEvent receives an event.Event proto definition that will be routed correctly.
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

//createAuthor creates an Author on the database.
func (rm *RouteManager) createAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	uuid, err := rm.connector.AddAuthor(database.AuthorFromGrpc(author))
	parsedUuid := uuid.String()
	return []string{parsedUuid}, err
}

// updateAuthor updates an Author on the database.
func (rm *RouteManager) updateAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	err := rm.connector.UpdateAuthor(database.AuthorFromGrpc(author))
	return nil, err
}

// readAuthor query one or all authors from the database.
func (rm *RouteManager) readAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.AllEntries {
		return rm.readAllAuthors()
	} else {
		return rm.readAuthorById(query.GetUuid())
	}
}

// readAuthorById get an Author entry by ID.
func (rm *RouteManager) readAuthorById(uuid string) ([]string, error) {
	author, err := rm.connector.GetAuthor(uuid)
	parsedAuthor := database.AuthorToGrpc(*author)
	return []string{utils.EncodeAuthorToString(parsedAuthor)}, err
}

// readAllAuthors retrieve all Author entries from the database.
func (rm *RouteManager) readAllAuthors() ([]string, error) {
	authors := rm.connector.GetAuthors()
	parsedAuthors := database.AuthorListToGrpcList(authors)
	return []string{utils.EncodeAuthorsListToString(&parsedAuthors)}, nil
}

// deleteAuthor deletes an Author from the database.
func (rm *RouteManager) deleteAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.Uuid == nil {
		return nil, errors.New("uuid not set on the request")
	}

	err := rm.connector.DeleteAuthor(query.GetUuid())
	return nil, err
}
