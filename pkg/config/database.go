package config

import "fmt"

type DBConfig struct {
	host          string
	port          int
	name          string
	user          string
	password      string
	migrationPath string
}

func (db *DBConfig) Address() string {
	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", db.host, db.user, db.password, db.name, db.port)
}

func (db *DBConfig) MigrationPath() string {
	return db.migrationPath
}

func newDBConfig() DBConfig {
	return DBConfig{
		host:          getString("DB_HOST", "localhost"),
		port:          getInt("DB_PORT", 5432),
		name:          getString("DB_NAME", "postgres"),
		user:          getString("DB_USER", "postgres"),
		password:      getString("DB_PASSWORD", "pwd"),
		migrationPath: getString("MIGRATION_PATH", ""),
	}
}
