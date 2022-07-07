package database

import (
	"errors"
	"log"
	"os"

	"marketplace/model"

	"github.com/jmoiron/sqlx"
)

type PostgresqlRepository struct {
	db string
}

func getConfig() (string, string, string, string, string, error) {
	log.Println("Fetching env vars for Database")

	host, hostPresent := os.LookupEnv("DATABASE_HOST")
	user, userPresent := os.LookupEnv("POSTGRES_USER")
	pass, passPresent := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbnamePresent := os.LookupEnv("POSTGRES_DB")
	port, portPresent := os.LookupEnv("DATABASE_PORT")

	if !hostPresent || !userPresent || !passPresent || !dbnamePresent || !portPresent {
		log.Println("Error occurred while fetching env vars")
		return "", "", "", "", "", errors.New("Error occurred while fetching env vars")
	}
	return host, user, pass, dbname, port, nil
}

func Connect() (*PostgresqlRepository, error) {
	log.Println("Connecting to DB")

	getConfig()

	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var schemas model.Schema
	createSchemas := schemas.GetCreateSchemas()

	// execute schema

	db.MustExec(createSchemas)

	return nil, nil
}

func (instance *PostgresqlRepository) Create(value interface{}, omits ...string) error {

	log.Println("Creating DB")
	return nil
}
