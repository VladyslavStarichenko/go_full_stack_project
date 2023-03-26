package main

import (
	"html/template"
	"net/http"
	"time"
)

type Performance struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"`
	Description string    `json:"description"`
}

type PerformanceInfo struct {
	Name        string `json:"performance_name"`
	Description string `json:"performance_description"`
	NumActors   int    `json:"num_actors"`
}

func getPerformanceStatistics() ([]*PerformanceInfo, error) {
	rows, err := database.Query(`SELECT p.name AS "Performance Name", p.description AS "Performance Description", COUNT(ap.actor_id) AS "Number of Actors" ` +
		`FROM performances p JOIN actor_performances ap ON p.id = ap.performance_id GROUP BY p.id ORDER BY "Number of Actors" DESC, "Performance Name" ASC `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	performanceStatistics := []*PerformanceInfo{}
	for rows.Next() {
		performanceStatistic := &PerformanceInfo{}
		err := rows.Scan(
			&performanceStatistic.Name,
			&performanceStatistic.Description,
			&performanceStatistic.NumActors,
		)
		if err != nil {
			return nil, err
		}
		performanceStatistics = append(performanceStatistics, performanceStatistic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performanceStatistics, nil
}

func performanceStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	performanceStatistics, err := getPerformanceStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	performanceStatistic := make([]PerformanceInfo, len(performanceStatistics))
	for i, statistic := range performanceStatistics {
		performanceStatistic[i] = *statistic
	}

	tmpl := template.Must(template.ParseFiles("templates/performance_statistics.html"))
	err = tmpl.Execute(w, performanceStatistic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
