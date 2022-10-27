package database

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"user-management/middleware"
	"user-management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresqlRepository struct {
	db *gorm.DB
}

func Connect() (*PostgresqlRepository, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Println("Error Connecting to the database: " + err.Error())
		return nil, err
	}

	log.Println("Connected to the Database")

	migrate(db, &models.User{})

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

func getConfig() (string, error) {
	log.Println("Fetching env vars for Database")

	host, hostPresent := os.LookupEnv("DATABASE_HOST")
	user, userPresent := os.LookupEnv("POSTGRES_USER")
	pass, passPresent := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbnamePresent := os.LookupEnv("POSTGRES_DB")
	port, portPresent := os.LookupEnv("DATABASE_PORT")

	if !hostPresent || !userPresent || !passPresent || !dbnamePresent || !portPresent {
		log.Println("Error occurred while fetching env vars")
		return "", middleware.NewError(http.StatusInternalServerError, "Error occurred while fetching env vars")
	}

	log.Println("host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port)
	return "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port, nil
}

func isSliceOrArray(value interface{}) bool {
	return reflect.ValueOf(value).Elem().Kind() == reflect.Slice || reflect.ValueOf(value).Elem().Kind() == reflect.Array
}

func (instance *PostgresqlRepository) Create(value interface{}) error {

	result := instance.db.Create(value)
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
			return middleware.NewError(http.StatusNotFound, "Record Not found")
		}
		return result.Error
	}

	log.Println("Fetched database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Update(value interface{}) error {
	result := instance.db.Save(value)
	if result.Error != nil {
		log.Println("Error while updating a database entry: " + fmt.Sprintf("%v", value))
		return result.Error
	}

	log.Println("Updated database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Delete(value interface{}) error {

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
