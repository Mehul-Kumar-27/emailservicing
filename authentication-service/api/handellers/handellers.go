package handellers

import (
	data "authentication-service/api/data"
	"database/sql"
)

type DatabaseModel struct {
	DB     *sql.DB
	Models data.Models
}

func NewDatabaseModel(db *sql.DB) *DatabaseModel {
	return &DatabaseModel{
		DB:     db,
		Models: data.New(db),
	}
}

func main() {
}
