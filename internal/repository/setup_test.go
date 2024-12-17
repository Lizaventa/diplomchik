package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	configuration "github.com/liza/labwork_45/internal/config"
)

const (
	pgUsername = "eugen"
	pgPassword = "ur2qly1ini"
	pgDB       = "lab"
)

// SetupTestPgx function to test pgx methods
func SetupTestPgx() (*pgxpool.Pool, func(), error) {
	dbAddr, err := configuration.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return nil, nil, fmt.Errorf("Error extracting env variables: %v", err)
	}
	// Initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig(dbAddr.PgxDBAddr)
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return nil, nil, fmt.Errorf("Error extracting env variables: %v", err)
	}

	// Establishing a new connection to a PostgreSQL database using the pgx driver
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	// Output to console
	fmt.Println("Connected to PostgreSQL!")

	cleanup := func() {
		dbpool.Close()
		//pool.Purge(resource)
	}

	return dbpool, cleanup, nil
}

// TestMain is an integration test for Application<->database conenction
func TestMain(m *testing.M) {
	dbpool, cleanupPgx, err := SetupTestPgx()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
		cleanupPgx()
		os.Exit(1)
	}
	rps = NewPsqlConnection(dbpool)
	exitVal := m.Run()
	cleanupPgx()
	os.Exit(exitVal)
}

func CreateTestProfile() (uuid.UUID, error) {
	id, err := rps.InsertUser(context.Background(), testProfile)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func DeleteTestProfile(id uuid.UUID) error {
	err := rps.DeleteUserByID(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}
