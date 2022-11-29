package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"myproject-page/connection"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer((http.Dir("./public")))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", formEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", editProject).Methods("POST")

	fmt.Println(("server berjalan di port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

type Project struct {
	ID                     int
	Title                  string
	DateStart              time.Time
	DateEnd                time.Time
	Format_date_start      string
	Format_date_start_edit string
	Format_date_end        string
	Description            string
}

// var projects = []Project{}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	dateStart := r.PostForm.Get("date-start")
	dateEnd := r.PostForm.Get("date-end")

	_, errQuery := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(title, start_date, end_date, description) VALUES ($1, $2, $3, $4)", title, dateStart, dateEnd, content)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// var newProject = Project{
	// 	DateStartEdit: dateStart,
	// 	DateEndEdit: dateEnd,
	// }

	// projects = append(projects, newProject)

	fmt.Println(dateStart)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/editProject.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	var ProjectEdit = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id = $1", index).Scan(&ProjectEdit.ID, &ProjectEdit.Title, &ProjectEdit.DateStart, &ProjectEdit.DateEnd, &ProjectEdit.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	ProjectEdit.Format_date_start = ProjectEdit.DateStart.Format("2 January 2006")
	ProjectEdit.Format_date_end = ProjectEdit.DateEnd.Format("2 January 2006")

	dataEdit := map[string]interface{}{
		"Project": ProjectEdit,
	}

	tmpt.Execute(w, dataEdit)
}

func editProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	dateStart := r.PostForm.Get("date-start")
	dateEnd := r.PostForm.Get("date-end")

	_, errQuery := connection.Conn.Exec(context.Background(),
		"UPDATE public.tb_projects SET title=$1, start_date=$2, end_date=$3, description=$4 WHERE id = $5", title, dateStart, dateEnd, content, index)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// var newProject = Project{
	// 	Title:       title,
	// 	Description: content,
	// 	DateStart:   dateStart,
	// 	DateEnd:     dateEnd,
	// }

	// // projects = append(projects, newProject)
	// projects[index] = newProject

	// fmt.Println(index)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects")
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.ID, &each.Title, &each.DateStart, &each.DateEnd, &each.Description)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

		each.Format_date_start = each.DateStart.Format("2 January 2006")
		each.Format_date_end = each.DateEnd.Format("2 January 2006")
		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Projects": result,
	}

	// dataProject := map[string]interface{}{
	// 	"Projects": projects,
	// }

	tmpt.Execute(w, resData)
}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id = $1", id).Scan(&ProjectDetail.ID, &ProjectDetail.Title, &ProjectDetail.DateStart, &ProjectDetail.DateEnd, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	ProjectDetail.Format_date_start = ProjectDetail.DateStart.Format("2 January 2006")
	ProjectDetail.Format_date_end = ProjectDetail.DateEnd.Format("2 January 2006")

	dataDetail := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpt.Execute(w, dataDetail)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", index)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}
