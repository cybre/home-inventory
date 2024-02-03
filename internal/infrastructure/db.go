package infrastructure

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewCassandraSession(hosts []string, keyspace string) (*gocql.Session, error) {
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

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "cassandra", databaseInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return session, nil
}
