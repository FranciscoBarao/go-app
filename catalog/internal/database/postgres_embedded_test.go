package database

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	embedded "github.com/fergusstrange/embedded-postgres"
)

// TestPostgres is an embedded postgres database for testing
type TestPostgres struct {
	Database *embedded.EmbeddedPostgres
	Port     int
}

// Shutdown stops the database instance and releases allocated resources
func (p *TestPostgres) Shutdown() {
	if p.Database == nil {
		return
	}

	_ = p.Database.Stop()
}

// NewTestPostgres constructs a new test database
func NewTestPostgres(uniqueID string) (*TestPostgres, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to find a port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	runtimeDirectory := filepath.Join(userHome, ".embedded-postgres-go", uniqueID)

	postgresConfig := embedded.
		DefaultConfig().
		Port(uint32(port)).
		RuntimePath(runtimeDirectory)
	db := embedded.NewDatabase(postgresConfig)
	if err := db.Start(); err != nil {
		return nil, fmt.Errorf("failed to initialise postgres: %w", err)
	}

	return &TestPostgres{
		Database: db,
		Port:     port,
	}, nil
}
