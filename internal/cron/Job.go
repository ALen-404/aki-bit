package cron

import (
	"btc_order/pkg/postgresql"
)

type Job struct {
	db *postgresql.Database
}

func New(db *postgresql.Database) *Job {
	return &Job{db: db}
}

func (j *Job) Start() {
	go func() {
		j.btcOrderJob()
	}()
}
