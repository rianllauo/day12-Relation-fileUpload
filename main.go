package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"myproject-page/connection"
	"myproject-page/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer((http.Dir("./public")))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/", middleware.UploadFile(addProject)).Methods("POST")
	// route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", formEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", editProject).Methods("POST")

	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println(("server berjalan di port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

type MetaData struct {
	Id        int
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	//card project struct
	ID                int
	Title             string
	Creator           string
	DateStart         time.Time
	DateEnd           time.Time
	Duration          string
	Month             float64
	Format_date_start string
	Format_date_end   string
	Description       string
	Image             string
	Technologies      []string
	NodeJs            string
	ReactJs           string
	NextJs            string
	Javascript        string
	IsLogin           bool
}

type Login struct {
	Er           string
	InputInvalid string
}

var DataErr = Login{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// sessions
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Names"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, title, start_date, end_date, description, technologies, image, tb_user.name as creator FROM tb_projects LEFT JOIN tb_user ON tb_projects.user_id = tb_user.id ORDER BY id DESC")
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.ID, &each.Title, &each.DateStart, &each.DateEnd, &each.Description, &each.Technologies, &each.Image, &each.Creator)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

		diff := each.DateEnd.Sub(each.DateStart)
		days := diff.Hours() / 24
		month := math.Floor(days / 30)

		dy := strconv.FormatFloat(days, 'f', 0, 64)
		mo := strconv.FormatFloat(month, 'f', 0, 64)

		if days < 30 {
			each.Duration = dy + " Days"
		} else if days > 30 {
			each.Duration = mo + " Month"
		}

		// checked condition
		if each.Technologies[0] == "nodejs" {
			each.NodeJs = "nodejs.svg"
		} else {
			each.NodeJs = "d-none"
		}
		if each.Technologies[1] == "nextjs" {
			each.NextJs = "nextjs.svg"
		} else {
			each.NextJs = "d-none"
		}
		if each.Technologies[2] == "react" {
			each.ReactJs = "react.svg"
		} else {
			each.ReactJs = "d-none"
		}
		if each.Technologies[3] == "javascript" {
			each.Javascript = "javascript.svg"
		} else {
			each.Javascript = "d-none"
		}

		each.Format_date_start = each.DateStart.Format("2 January 2006")
		each.Format_date_end = each.DateEnd.Format("2 January 2006")

		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Projects": result,
		"Data":     Data,
	}
	tmpt.Execute(w, resData)

}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// sessions
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Names"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	Data := map[string]interface{}{
		"DataFlash": Data,
		// "DataFlash": DataFlash,
	}

	tmpt.Execute(w, Data)
}

// SELECT name FROM tb_user
// LEFT JOIN tb_projects ON tb_user.id = tb_projects.id;

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	dateStart := r.PostForm.Get("date-start")
	dateEnd := r.PostForm.Get("date-end")

	dataContext := r.Context().Value("dataImages")
	image := dataContext.(string)

	nodeJs := r.PostForm.Get("nodeJs")
	nextJs := r.PostForm.Get("nextJs")
	reactJs := r.PostForm.Get("reactJs")
	javascript := r.PostForm.Get("javascript")

	checked := []string{
		nodeJs,
		nextJs,
		reactJs,
		javascript,
	}

	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	userPost := session.Values["Id"]

	fmt.Println(userPost)

	_, errQuery := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(title, start_date, end_date, description, technologies, image, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", title, dateStart, dateEnd, content, checked, image, userPost)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

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

	// ProjectEdit.nodeJs = ProjectEdit.Technologies[0]

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

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var each = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, technologies FROM tb_projects WHERE id = $1", id).Scan(&each.ID, &each.Title, &each.DateStart, &each.DateEnd, &each.Description, &each.Technologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	diff := each.DateEnd.Sub(each.DateStart)
	days := diff.Hours() / 24
	month := math.Floor(days / 30)

	dy := strconv.FormatFloat(days, 'f', 0, 64)
	mo := strconv.FormatFloat(month, 'f', 0, 64)

	if days < 30 {
		each.Duration = dy + " Days"
	} else if days > 30 {
		each.Duration = mo + " Month"
	}

	// checked condition
	if each.Technologies[0] == "nodejs" {
		each.NodeJs = "nodejs.svg"
	} else {
		each.NodeJs = "d-none"
	}
	if each.Technologies[1] == "nextjs" {
		each.NextJs = "nextjs.svg"
	} else {
		each.NextJs = "d-none"
	}
	if each.Technologies[2] == "react" {
		each.ReactJs = "react.svg"
	} else {
		each.ReactJs = "d-none"
	}
	if each.Technologies[3] == "javascript" {
		each.Javascript = "javascript.svg"
	} else {
		each.Javascript = "d-none"
	}

	fmt.Println(diff)

	each.Format_date_start = each.DateStart.Format("2 January 2006")
	each.Format_date_end = each.DateEnd.Format("2 January 2006")

	dataDetail := map[string]interface{}{
		"Project": each,
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

	// sessions
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Names"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	Data := map[string]interface{}{
		"DataFlash": Data,
		// "DataFlash": DataFlash,
	}

	tmpt.Execute(w, Data)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/register.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
	}

	tmpt.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")

	password := r.PostForm.Get("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(),
		"INSERT INTO tb_user(name, email, password) VALUES($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	session.AddFlash("successfully registered!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
	}

	tmpt.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email = $1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		DataErr.Er = "Password Wrong"
		DataErr.InputInvalid = "is-invalid"
		t := template.Must(template.ParseFiles("views/login.html"))
		Data := map[string]interface{}{
			"DataErr": DataErr,
		}
		t.Execute(w, Data)
	}

	session.Values["IsLogin"] = true
	session.Values["Names"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 30

	session.AddFlash("Successfully login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
