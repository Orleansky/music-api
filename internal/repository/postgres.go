package repository

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewStorage(cfg Config) (*pgxpool.Pool, error) {
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	logrus.WithField("connstr", connstr).Debug("Connecting to the database")

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to the database")
		return nil, err
	}

	logrus.Info("Connected to the database")

	err = migration(db)
	if err != nil {
		logrus.WithError(err).Error("Failed to run migrations")
		return nil, err
	}

	logrus.Info("Migrations completed successfully")

	return db, nil
}

func migration(db *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*db.Config().ConnConfig)
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logrus.WithError(err).Error("Failed to create Postgres driver")
		return err
	}

	logrus.Debug("Postgres driver created successfully")

	sourceDriver, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		logrus.WithError(err).Error("Failed to open migration files")
		return err
	}

	logrus.Debug("Migration files opened successfully")

	m, err := migrate.NewWithInstance(
		"postgres",
		sourceDriver,
		"postgres",
		driver,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to create migration instance")
		return err
	}

	logrus.Debug("Migration instance created successfully")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.WithError(err).Error("Failed to apply migrations")
		return err
	}

	logrus.Info("Migrations applied successfully")

	return nil
}
