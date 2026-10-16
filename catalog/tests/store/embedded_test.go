package store_test

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	embedded "github.com/fergusstrange/embedded-postgres"
)

func startTestPostgres(uniqueID string) (int, func(), error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to find a port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	userHome, err := os.UserHomeDir()
	if err != nil {
		return 0, nil, err
	}
	runtimeDirectory := filepath.Join(userHome, ".embedded-postgres-go", uniqueID)

	postgresConfig := embedded.
		DefaultConfig().
		Port(uint32(port)).
		RuntimePath(runtimeDirectory)
	database := embedded.NewDatabase(postgresConfig)
	if err := database.Start(); err != nil {
		return 0, nil, fmt.Errorf("failed to initialise postgres: %w", err)
	}

	return port, func() { _ = database.Stop() }, nil
}
