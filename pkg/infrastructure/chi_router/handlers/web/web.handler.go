package web

import (
	"fmt"
	"go-clean-api/utils"
	"html/template"
	"net/http"

	"github.com/go-chi/render"
)

// HealthCheck returns status code 200
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// GetAPIv1Doc returns the API v1 documentation
func GetAPIv1Doc(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./templates/doc_api_v1.gohtml")
	if err != nil {
		return utils.Err500(w, err, "Error when parsing HTML template", nil)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		return utils.Err500(w, err, "Error when executing HTML template", nil)
	}

	return nil
}

// BigTasks returns a big JSON
// TODO: Remove?
func BigTasks(w http.ResponseWriter, r *http.Request) error {
	type Task struct {
		ID   int
		Name string
	}

	var tasks []Task
	for i := 0; i < 10_000; i++ {
		tasks = append(tasks, Task{
			ID:   i*100_000 + 1,
			Name: fmt.Sprintf("My task with ID: %d", i),
		})
	}

	render.JSON(w, r, tasks)

	return nil
}
