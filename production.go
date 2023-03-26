package main

import "time"

type Production struct {
	ID                  int64              `json:"id"`
	PerformanceID       int64              `json:"performance_id"`
	ProductionCompanyID int64              `json:"production_company_id"`
	StartDate           time.Time          `json:"start_date"`
	EndDate             *time.Time         `json:"end_date,omitempty"`
	Budget              int                `json:"budget"`
	Performance         *Performance       `json:"performance,omitempty"`
	ProductionCompany   *ProductionCompany `json:"production_company,omitempty"`
}
