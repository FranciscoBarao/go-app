package database

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"marketplace/middleware"
	"marketplace/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresqlRepository struct {
	db *sqlx.DB
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

	return "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " port=" + port + " sslmode=disable", nil
}

func Connect() (*PostgresqlRepository, error) {
	log.Println("Connecting to DB")

	config, err := getConfig()
	if err != nil {
		log.Println("Error Getting environment variables: " + err.Error())
		return nil, err
	}

	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Println("Error Connecting to the database: " + err.Error())
		return nil, err
	}

	log.Println("Connected to the Database")

	var schemas model.Schema = model.SchemaAgregator{}
	createSchemas := schemas.GetCreateSchemas()

	_, err = db.Exec(createSchemas)
	if err != nil {
		log.Println("Error Creating schemas " + err.Error())
		return nil, err
	}

	log.Println("Database Schemas creation Completed")

	return &PostgresqlRepository{db}, nil
}

func (instance *PostgresqlRepository) Create(query string, value ...interface{}) error {

	_, err := instance.db.Exec(query, value...)
	if err != nil {
		log.Println("Error while creating a database entry: " + fmt.Sprintf("%v", query))
		log.Println(err)
		return err
	}

	log.Println("Created database entry: " + fmt.Sprintf("%v", value))
	return nil
}

func (instance *PostgresqlRepository) ReadAll(query string, value interface{}) error {

	err := instance.db.Select(value, query)
	if err != nil {
		log.Println("Error fetching database entries: " + fmt.Sprintf("%v", query))
		log.Println(err)
		return err
	}
	log.Println("Fetched database entry: " + fmt.Sprintf("%v", value))
	return nil
}
