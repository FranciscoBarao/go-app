package store_test

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	embedded "github.com/fergusstrange/embedded-postgres"
)

type testPostgres struct {
	database *embedded.EmbeddedPostgres
	port     int
}

func (p *testPostgres) Shutdown() {
	if p == nil || p.database == nil {
		return
	}
	_ = p.database.Stop()
}

func newTestPostgres(uniqueID string) (*testPostgres, error) {
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

	return &testPostgres{
		database: db,
		port:     port,
	}, nil
}
