package database

import (
	"errors"
	"fmt"
	"go-app/model"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	migrate(db, model.BoardGame{})

	log.Println("Database Migration Completed")

	return &PostgresqlRepository{db}, nil
}

func migrate(db *gorm.DB, value interface{}) error {
	err := db.AutoMigrate(&value)
	if err != nil {
		log.Println("Error migrating database: " + fmt.Sprintf("%v", value))
		return err
	}

	log.Println("Migrated " + fmt.Sprintf("%v", value))
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
		return "", errors.New("Error occurred while fetching env vars")
	}

	return "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port, nil
}

func (instance *PostgresqlRepository) Create(value interface{}) error {

	result := instance.db.Create(value) //value and not &value ->  You want to pass a pointer of the struct not the interface
	if result.Error != nil {
		log.Println("Error while creating a database entry: " + fmt.Sprintf("%v", value))
		return result.Error
	}

	log.Println("Created database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) Read(value interface{}, search string, identifier string) error {
	result := instance.db.First(&value, search, identifier)
	if result.Error != nil {
		log.Println("Error while fetching a database entry: " + search + " " + identifier)
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
	result := instance.db.Delete(value)
	if result.Error != nil {
		log.Println("Error while deleting a database entry: " + fmt.Sprintf("%v", value))
		return result.Error
	}

	log.Println("Deleted database entry: " + fmt.Sprintf("%v", value))
	return nil
}
