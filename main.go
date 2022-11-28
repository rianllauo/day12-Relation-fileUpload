package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"myproject-page/connection"
	"net/http"
	"strconv"

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
	Id          int
	Title       string
	DateStart   string
	DateEnd     string
	Description string
}

var projects = []Project{
	{
		Title:       "Aplikasi web dumbways",
		DateStart:   "11 november 2022",
		DateEnd:     "12 desember 2022",
		Description: "lorem ipsum dolor si amet",
		// NodeJs:        "public/img/nodejs.svg",
	},
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// dateStart := r.PostForm.Get("date-start")
	// dateEnd := r.PostForm.Get("date-end")

	// var newProject = Project{
	// 	Title:       title,
	// 	Description: content,
	// 	DateStart:   dateStart,
	// 	DateEnd:     dateEnd,
	// }

	dataProject, errQuery := connection.Conn.Query(context.Background(), `INSERT INTO public.tb_projects(
		id, title, description)
		VALUES ( 'title coba', 'lorem ipsum dolor si amet' )`)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}
	// projects = append(projects, newProject)

	fmt.Println(dataProject)

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

	for i, data := range projects {
		if i == index {
			ProjectEdit = Project{
				Id:          i,
				Title:       data.Title,
				Description: data.Description,
				DateStart:   data.DateStart,
				DateEnd:     data.DateEnd,
			}
		}
	}

	// fmt.Println(projects)

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

	var newProject = Project{
		Title:       title,
		Description: content,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
	}

	// projects = append(projects, newProject)
	projects[index] = newProject

	fmt.Println(index)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title, description FROM tb_projects")
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Id, &each.Title, &each.Description)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

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

	for index, data := range projects {
		if index == id {
			ProjectDetail = Project{
				Title:       data.Title,
				Description: data.Description,
				DateStart:   data.DateStart,
				DateEnd:     data.DateEnd,
			}
		}
	}

	dataDetail := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpt.Execute(w, dataDetail)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	projects = append(projects[:index], projects[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
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
