package caskin

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBOption holds the configuration for the metadata database connection.
// Pass it as part of [Options] when constructing a caskin service via [New].
type DBOption struct {
	// DSN is the data source name (connection string) for the database.
	DSN string `json:"dsn"`
	// Type selects the database driver. Supported values: "sqlite", "mysql"
	// (default when empty), and "postgres".
	Type string `json:"type"`
}

// NewDB opens a GORM database connection using the configured driver and DSN.
func (o *DBOption) NewDB() (*gorm.DB, error) {
	dialect, err := getDialect(o.Type, o.DSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getDialect(ty string, dsn string) (gorm.Dialector, error) {
	switch ty {
	case "sqlite":
		return sqlite.Open(dsn), nil
	case "mysql", "":
		return mysql.Open(dsn), nil
	case "postgres":
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %v", ty)
	}
}
