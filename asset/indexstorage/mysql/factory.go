package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci-inspector/common/mysql"
)

const (
	// language=MySQL
	createTableSQL = `
CREATE TABLE IF NOT EXISTS job_assets (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	job_id BIGINT NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
	UNIQUE u_asset(job_id, asset_name),
	INDEX i_job_id(job_id)
)
`
)

// NewMySQLAssetIndex creates a MySQL storage for asset indexes.
func NewMySQLAssetIndex(config mysql.Config, logger *log.Logger) (indexstorage.AssetIndex, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open(
		"mysql",
		config.ConnectString(),
	)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(
		`CREATE DATABASE IF NOT EXISTS ` + config.Database,
	); err != nil {
		return nil, fmt.Errorf("failed to create database (%w)", err)
	}
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create assets table (%w)", err)
	}

	return &mysqlAssetIndex{
		db:     db,
		logger: logger,
	}, nil
}