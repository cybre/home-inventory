package cassandra

import (
	"errors"
	"fmt"
	"os"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewSession(hosts []string, keyspace string) (*gocql.Session, error) {
	if err := createKeyspace(hosts, keyspace); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	databaseInstance, err := cassandra.WithInstance(session, &cassandra.Config{
		MigrationsTable:       "schema_migrations",
		KeyspaceName:          keyspace,
		MultiStatementEnabled: true,
		MultiStatementMaxSize: cassandra.DefaultMultiStatementMaxSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database instance for migrations: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://migrations/%s", keyspace), "cassandra", databaseInstance)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return session, nil
		}

		return nil, fmt.Errorf("failed to create migrations: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return session, nil
}

func createKeyspace(hosts []string, keyspace string) error {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "system"

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	if err := session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1}", keyspace)).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	return nil
}
