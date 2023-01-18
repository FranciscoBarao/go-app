package database

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"rating-service/middleware"
	"rating-service/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresqlRepository struct {
	db *gorm.DB
}

func Connect(isTest bool) (*PostgresqlRepository, error) {
	config, err := getConfig(isTest)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Println("Error Connecting to the database: " + err.Error())
		return nil, err
	}

	log.Println("Connected to the Database")

	migrate(db, &model.Rating{})

	log.Println("Database Migration Completed")

	return &PostgresqlRepository{db}, nil
}

func migrate(db *gorm.DB, model interface{}) error {
	if err := db.AutoMigrate(model); err != nil {
		log.Println("Error migrating database: " + fmt.Sprintf("%v", model))
		return err
	}

	log.Println("Migrated " + fmt.Sprintf("%v", model))
	return nil
}

func getConfig(isTest bool) (string, error) {
	log.Println("Fetching env vars for Database")

	host, hostPresent := os.LookupEnv("DATABASE_HOST")
	user, userPresent := os.LookupEnv("POSTGRES_USER")
	pass, passPresent := os.LookupEnv("POSTGRES_PASSWORD")
	port, portPresent := os.LookupEnv("DATABASE_PORT")
	dbname, dbnamePresent := os.LookupEnv("POSTGRES_DB")

	if isTest {
		dbname, dbnamePresent = "test", true
	}

	if !hostPresent || !userPresent || !passPresent || !dbnamePresent || !portPresent {
		log.Println("Error occurred while fetching env vars")
		return "", middleware.NewError(http.StatusInternalServerError, "Error occurred while fetching env vars")
	}

	return "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port, nil
}

func isSliceOrArray(value interface{}) bool {
	return reflect.ValueOf(value).Elem().Kind() == reflect.Slice || reflect.ValueOf(value).Elem().Kind() == reflect.Array
}

func (instance *PostgresqlRepository) Create(value interface{}, omits ...string) error {

	result := instance.db.Omit(omits...).Create(value)
	if result.Error != nil {
		log.Println("Error while creating a database entry: " + fmt.Sprintf("%v", value))
		if errors.Is(result.Error, gorm.ErrRegistered) {
			return middleware.NewError(http.StatusConflict, "Entry already registered")
		}
		return result.Error
	}

	log.Println("Created database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Read(value interface{}, sort, search, identifier string) error {

	var result *gorm.DB
	if isSliceOrArray(value) {
		if search == "" {
			result = instance.db.Preload(clause.Associations).Order(sort).Find(value) // Find all with sort and NO filters
		} else {
			result = instance.db.Preload(clause.Associations).Order(sort).Find(value, search, identifier) // Find all with filters and sort
		}
	} else {
		result = instance.db.Preload(clause.Associations).First(value, search, identifier) // Find 1 Specific
	}

	if result.Error != nil {
		log.Println("Error while reading a database entry: " + search + " " + identifier)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("Error Record not found: " + search + " " + identifier)
			return middleware.NewError(http.StatusNotFound, "Record not found")
		}
		return result.Error
	}

	log.Println("Fetched database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Update(value interface{}, omits ...string) error {
	result := instance.db.Omit(omits...).Save(value)
	if result.Error != nil {
		log.Println("Error while updating a database entry: " + fmt.Sprintf("%v", value))
		return result.Error
	}

	log.Println("Updated database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Delete(value interface{}) error {

	// Delete BG and all its associations (E.g Tags associations)
	result := instance.db.Select(clause.Associations).Delete(value)
	if result.Error != nil {
		log.Println("Error while deleting a database entry: " + fmt.Sprintf("%v", value))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return middleware.NewError(http.StatusNotFound, "Record Not found")
		}
		return result.Error
	}

	log.Println("Deleted database entry: " + fmt.Sprintf("%v", value))
	return nil
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<        ASSOCIATIONS        >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// The following section presents the associations generic methods. This section is up for debate and will possibly change in the future.

// Method that Adds certain associations to a certain model (E.g Add Tags to a Boardgame)
func (instance *PostgresqlRepository) AppendAssociatons(model interface{}, association string, values interface{}) error {

	err := instance.db.Model(model).Association(association).Append(values)
	if err != nil {
		log.Println("Error while appending associations of type: " + association + " to model: " + fmt.Sprintf("%v", model) + " with values: " + fmt.Sprintf("%v", values))
		log.Println(err)
		return err
	}

	log.Println("Associated: " + association + " to model: " + fmt.Sprintf("%v", model) + " with values: " + fmt.Sprintf("%v", values))
	return nil
}

// Method that Gets associations of a type of a certain model (E.g Get Tags of a Boardgame)
func (instance *PostgresqlRepository) ReadAssociatons(model interface{}, association string, store interface{}) error {

	err := instance.db.Model(model).Association(association).Find(store)
	if err != nil {
		log.Println("Error while Reading associations of type: " + association + " of model: " + fmt.Sprintf("%v", model))
		return err
	}

	log.Println("Fetched: " + association + " og model: " + fmt.Sprintf("%v", model) + " with values: " + fmt.Sprintf("%v", store))
	return nil
}

// Method that Replaces the values of a certain association of a certain model (E.g Replace Tags of a Boardgame)
func (instance *PostgresqlRepository) ReplaceAssociatons(model interface{}, association string, values interface{}) error {

	err := instance.db.Model(model).Association(association).Replace(values)
	if err != nil {
		log.Println("Error while replacing associations type: " + association + " from model: " + fmt.Sprintf("%v", model) + " with values: " + fmt.Sprintf("%v", values))
		return err
	}

	log.Println("Associated: " + association + " to model: " + fmt.Sprintf("%v", model) + " with values: " + fmt.Sprintf("%v", values))
	return nil
}

// Method that Deletes all values of a certain association of a certain model (E.g Delete all Tags of a Boardgame)
func (instance *PostgresqlRepository) DeleteAssociatons(model interface{}, association string) error {

	err := instance.db.Model(model).Association(association).Clear()
	if err != nil {
		log.Println("Error while deleting associations type: " + association + " from model: " + fmt.Sprintf("%v", model))
		return err
	}

	log.Println("Deleted Associations: " + association + " to model: " + fmt.Sprintf("%v", model))
	return nil
}
