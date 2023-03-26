package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Actor struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	Age                 int        `json:"age"`
	Gender              string     `json:"gender"`
	Nationality         string     `json:"nationality"`
	ContactPhone        string     `json:"contact_phone"`
	ContactEmail        string     `json:"contact_email"`
	ContactAddress      string     `json:"contact_address"`
	EmploymentStartDate time.Time  `json:"employment_start_date"`
	EmploymentEndDate   *time.Time `json:"employment_end_date,omitempty"`
	EmploymentSalary    int64      `json:"employment_salary"`
}

func DeleteActorHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the actor to delete from the URL path parameter
	idStr := mux.Vars(r)["id"]

	// Parse the ID string to an int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Println(id)
	// Delete the actor from the database
	_, err = database.Exec("DELETE FROM actors WHERE id = $1", id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Redirect the user back to the actors list page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditActorPage(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	actorID := vars["id"]

	row := database.QueryRow("SELECT id, name, age, gender, nationality, contact_phone, contact_email, contact_address, employment_salary FROM actors WHERE id = $1", actorID)

	actor := Actor{}

	err := row.Scan(
		&actor.ID,
		&actor.Name,
		&actor.Age,
		&actor.Gender,
		&actor.Nationality,
		&actor.ContactPhone,
		&actor.ContactEmail,
		&actor.ContactAddress,
		&actor.EmploymentSalary,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit_actor.html")
		tmpl.Execute(w, actor)
	}
}

func EditActorHandler(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		log.Println(err)
	}

	actorID := request.FormValue("id")

	_, err = database.Exec("UPDATE actors SET name=$1, age=$2, gender=$3, nationality=$4, contact_phone=$5, contact_email=$6, contact_address=$7, employment_salary=$8 WHERE id=$9",
		request.FormValue("name"),
		request.FormValue("age"),
		request.FormValue("gender"),
		request.FormValue("nationality"),
		request.FormValue("contact_phone"),
		request.FormValue("contact_email"),
		request.FormValue("contact_address"),
		request.FormValue("employment_salary"),
		actorID,
	)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, request, "/", 301)
}

func createActorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Extract actor data from form
		name := r.FormValue("name")
		age := r.FormValue("age")
		gender := r.FormValue("gender")
		nationality := r.FormValue("nationality")
		contactPhone := r.FormValue("contact_phone")
		contactEmail := r.FormValue("contact_email")
		contactAddress := r.FormValue("contact_address")
		employmentSalary := r.FormValue("employment_salary")

		fmt.Println(name, age, gender, nationality, contactPhone, contactEmail, contactAddress, employmentSalary)
		_, err = database.Exec(
			"INSERT INTO actors (name, age, gender, nationality, contact_phone, contact_email, contact_address, employment_salary)"+
				" VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
			name, age, gender, nationality, contactPhone, contactEmail, contactAddress, employmentSalary)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
		// Redirect to actors page
	} else {
		http.ServeFile(w, r, "templates/create_actor.html")

	}
}

func actorsHandler(w http.ResponseWriter, r *http.Request) {
	actorsPtrs, err := getUserData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actors := make([]Actor, len(actorsPtrs))
	for i, actorPtr := range actorsPtrs {
		actors[i] = *actorPtr
	}

	tmpl := template.Must(template.ParseFiles("templates/actors.html"))
	err = tmpl.Execute(w, actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUserData() ([]*Actor, error) {

	rows, err := database.Query("SELECT * FROM actors ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := []*Actor{}
	for rows.Next() {
		actor := &Actor{}
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.Age,
			&actor.Gender,
			&actor.Nationality,
			&actor.ContactPhone,
			&actor.ContactEmail,
			&actor.ContactAddress,
			&actor.EmploymentStartDate,
			&actor.EmploymentEndDate,
			&actor.EmploymentSalary,
		)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}
