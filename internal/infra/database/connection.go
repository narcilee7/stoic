package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds the database configuration
type Config struct {
	DSN         string `mapstructure:"dsn"`
	LogLevel    string `mapstructure:"log_level"` // debug, info, warn, error, silent
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

// DB wraps the gorm.DB type
type DB struct {
	*gorm.DB
}

// NewDB creates a new database connection using GORM
func NewDB(cfg Config) (*DB, error) {
	// Configure GORM logger
	gormConfig := &gorm.Config{
		PrepareStmt: true,
	}

	// Set log level
	switch cfg.LogLevel {
	case "debug":
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	case "info":
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	case "warn":
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	case "error":
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	default:
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(cfg.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign key constraints
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings (optional, adjust based on your needs)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)

	// Test the connection
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if err = db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign keys
	if err = db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable synchronous mode for better durability (slower but safer)
	if err = db.Exec("PRAGMA synchronous = NORMAL;").Error; err != nil {
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	fmt.Println("✅ Successfully connected to SQLite database with GORM")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		// Run PRAGMA optimize before closing
		_, _ = sqlDB.Exec("PREPARE optimize;")
		return sqlDB.Close()
	}
	return nil
}

// Migrate runs auto-migration for given models
func (db *DB) Migrate(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}

	// Enable foreign keys for migration
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations within a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			if err := tx.AutoMigrate(model); err != nil {
				return fmt.Errorf("failed to migrate model: %w", err)
			}
		}
		return nil
	})
}
