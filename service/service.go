package service

import (
	"github.com/go-pg/pg"

	"github.com/x-yield/over-api/tools"
)

type OverloadService struct {
	Db       *pg.DB
	InfluxDB *tools.InfluxDB
	S3       *tools.S3service
}

func NewOverloadService(db *pg.DB, influxdb *tools.InfluxDB, s3 *tools.S3service) *OverloadService {
	return &OverloadService{
		Db:       db,
		InfluxDB: influxdb,
		S3:       s3,
	}
}
