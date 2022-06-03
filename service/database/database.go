package database

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DbConnector connector used on the service.
type DbConnector struct {
	Database *gorm.DB
}

// Author type that will be stored in the DbConnector.
type Author struct {
	gorm.Model
	UUID   uuid.UUID `gorm:"primaryKey,default:uuid_generate_v4()"`
	Name   string
	PicURL *string
}

// NewDatabase Creates a new in memory DbConnector and automatically migrates the
// Author model.
func NewDatabase() DbConnector {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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

// CloseDatabase Closes that database that was open when creating a new database using the
// NewDatabase method.
func (database *DbConnector) CloseDatabase() {
	db, _ := database.Database.DB()
	defer db.Close()
}

// AddAuthor Adds an author to the database.
func (database *DbConnector) AddAuthor(id *string, name string, picUrl *string) (*uuid.UUID, error) {
	var _uuid uuid.UUID
	if id != nil {
		parsedUUID, err := uuid.Parse(*id)
		if err != nil {
			return nil, err
		}
		_uuid = parsedUUID
	} else {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		_uuid = newUUID
	}
	var author = Author{UUID: _uuid, Name: name, PicURL: picUrl}
	result := database.Database.Create(&author)
	return &author.UUID, result.Error
}

// GetAuthor Queries an author on the database using the uuid and return it to the caller.
func (database *DbConnector) GetAuthor(uuid string) (*Author, error) {
	var author *Author
	err := database.Database.First(&author, "uuid = ?", uuid).Error
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
func (database *DbConnector) UpdateAuthor(uuid string, name string, picUrl *string) error {
	var author, err = database.GetAuthor(uuid)
	if err != nil || author == nil {
		return err
	}
	err = database.Database.Model(author).Updates(Author{Name: name, PicURL: picUrl}).Error
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
