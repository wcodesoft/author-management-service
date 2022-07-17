package router

import (
	"errors"
	"service/database"
	"service/utils"

	eventProto "github.com/wcodesoft/event-manager/protos/go/event-manager.proto"
)

// RouteManager Object holding the necessary properties of the route manager.
type RouteManager struct {
	connector database.DbConnector
}

// NewRouteManager Creates a new RouteManager instance based on passed gorm DbConnector
func NewRouteManager(connector database.DbConnector) *RouteManager {
	return &RouteManager{
		connector: connector,
	}
}

// RouteEvent Process a received event from the message broker.
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

// createAuthor Creates an author from the information passed on the event.
func (rm *RouteManager) createAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	uuid, err := rm.connector.AddAuthor(database.AuthorFromGrpc(author))
	parsedUUID := uuid.String()
	return []string{parsedUUID}, err
}

// updateAuthor Updates an author with the new data passed on the event.
func (rm *RouteManager) updateAuthor(event *eventProto.Event) ([]string, error) {
	author := utils.DecodeAuthor(event.Message)
	err := rm.connector.UpdateAuthor(database.AuthorFromGrpc(author))
	return nil, err
}

// readAuthor Reads one or all authors from the database.
func (rm *RouteManager) readAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.AllEntries {
		return rm.readAllAuthors()
	}
	return rm.readAuthorByID(query.GetUuid())
}

// readAuthorByID Reads an author by the passed ID.
func (rm *RouteManager) readAuthorByID(uuid string) ([]string, error) {
	author, err := rm.connector.GetAuthor(uuid)
	parsedAuthor := database.AuthorToGrpc(*author)
	return []string{utils.EncodeAuthorToString(parsedAuthor)}, err
}

// readAllAuthors Reads all authors from the database.
func (rm *RouteManager) readAllAuthors() ([]string, error) {
	authors := rm.connector.GetAuthors()
	parsedAuthors := database.AuthorListToGrpcList(authors)
	return []string{utils.EncodeAuthorsListToString(&parsedAuthors)}, nil
}

// deleteAuthor Deletes one author from the database in case of a valid ID.
func (rm *RouteManager) deleteAuthor(event *eventProto.Event) ([]string, error) {
	query := utils.DecodeQuery(event.Message)
	if query.Uuid == nil {
		return nil, errors.New("uuid not set on the request")
	}

	err := rm.connector.DeleteAuthor(query.GetUuid())
	return nil, err
}
