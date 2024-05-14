package config

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	Database string
}

func NewPostgresConfig() (*PostgresConfig, error) {
	host, hostPresent := os.LookupEnv("DATABASE_HOST")
	user, userPresent := os.LookupEnv("POSTGRES_USER")
	pass, passPresent := os.LookupEnv("POSTGRES_PASSWORD")
	db, dbnamePresent := os.LookupEnv("POSTGRES_DB")
	port, portPresent := os.LookupEnv("DATABASE_PORT")

	if !hostPresent || !userPresent || !passPresent || !dbnamePresent || !portPresent {
		return nil, fmt.Errorf("failed to fetch postgres env vars")
	}

	return &PostgresConfig{
		Host:     host,
		Username: user,
		Password: pass,
		Database: db,
		Port:     port,
	}, nil
}

func (p *PostgresConfig) String() string {
	return "host=" + p.Host + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.Database + " port=" + p.Port
}
