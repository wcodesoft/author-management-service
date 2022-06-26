package database

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DbConnector connector used on the service.
type DbConnector struct {
	Database *gorm.DB
}

// Author type that will be stored in the DbConnector.
type Author struct {
	ID     *uuid.UUID `gorm:"primaryKey,unique,default:uuid_generate_v4()"`
	Name   string
	PicURL *string
}

// NewConnection Creates a new in memory DbConnector and automatically migrates the
// Author model.
func NewConnection(connector gorm.Dialector) DbConnector {
	db, err := gorm.Open(connector, &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}
	err = db.AutoMigrate(&Author{})
	if err != nil {
		panic("Failed to migrate to database.")
	}
	return DbConnector{
		Database: db,
	}
}

// uuidParseOrCreate Parse the string ID into a UUID or creates a new one when the passed value
// is invalid.
func uuidParseOrCreate(id string) uuid.UUID {
	value, err := uuid.Parse(id)
	if err == nil {
		return value
	}
	return uuid.New()
}

// CloseDatabase Closes that database that was open when creating a new database using the
// NewConnection method.
func (database *DbConnector) CloseDatabase() {
	db, _ := database.Database.DB()
	defer db.Close()
}

// AddAuthor Adds an author to the database.
func (database *DbConnector) AddAuthor(author Author) (*uuid.UUID, error) {
	authorToAdd := author
	if author.ID == nil {
		newUuid := uuid.New()
		author.ID = &newUuid
		authorToAdd = author
	}
	result := database.Database.Create(&authorToAdd)
	return author.ID, result.Error
}

// GetAuthor Queries an author on the database using the uuid and return it to the caller.
func (database *DbConnector) GetAuthor(uuid string) (*Author, error) {
	var author *Author
	err := database.Database.First(&author, "id = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return author, nil
}

// GetAuthors Gets all authors on the database.
func (database *DbConnector) GetAuthors() []Author {
	var allAuthors []Author
	database.Database.Find(&allAuthors)
	return allAuthors
}

// UpdateAuthor Updates the author entry with the new name and picUrl.
func (database *DbConnector) UpdateAuthor(author Author) error {
	if author.ID == nil {
		return errors.New("can´t update author without proper id")
	}
	var found, err = database.GetAuthor(author.ID.String())
	if err != nil || found == nil {
		return err
	}
	err = database.Database.Model(author).Updates(author).Error
	return err
}

// DeleteAuthor Deletes an author from the database with registered to the passed uuid.
func (database *DbConnector) DeleteAuthor(uuid string) error {
	var author, err = database.GetAuthor(uuid)
	if err != nil || author == nil {
		return err
	}
	err = database.Database.Delete(&author).Error
	return err
}
