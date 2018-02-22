package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// copied to dbcontroller
const (
	//DB_USER     = "postgres"
	//DB_PASSWORD = "111111"
	//DB_NAME     = "test"
	ServerPort   = "8081"
	FrontEndPort = "3000"
	IsDEV        = true
	ClientURL    = "http://localhost:3000"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		log.Panic("save error: ", err)
		return
	}
	renderTemplate(w, "view", p)
}

func hobbyHandler(w http.ResponseWriter, r *http.Request) {
	hobby := r.URL.Path[len("/hobby/"):]
	fmt.Fprintf(w, "Hi there, I love %s!", hobby)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
		log.Fatal("load error: ", err)
	}
	renderTemplate(w, "edit", p)

}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("save error: ", err)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "We are working on this page. Come back later!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err == nil {

			usrName := r.PostForm["username"][0]
			passWord := r.PostForm["password"][0]

			fmt.Println("username:", usrName)
			fmt.Println("password:", passWord)
			// TODO: implement login func

			if LoginPW(usrName, passWord) == true {
				//GenerateNewCookie(w, "uid", usrName)
				uData := make(map[string]string)
				uData["uid"] = usrName
				EncodeUserDataToCipherCookie(w, uData)

				scheme, hostname := GetRootURL(r)

				fmt.Println("r.URL")
				fmt.Println(r.URL.String())
				fmt.Println("r.host")
				fmt.Println(r.Host)
				fmt.Println("r.URL.Hostname")
				fmt.Println(r.URL.Hostname())
				redirectURL := scheme + "://" + hostname + ":" + FrontEndPort + "/"
				http.Redirect(w, r, redirectURL, http.StatusFound)

			} else {
				fmt.Println("Failed to login, check username/password")
			}
		} else {
			fmt.Println("Failed to parse form", err)
		}
	}
}

func GetRootURL(r *http.Request) (string, string) {
	hostname := r.URL.Hostname()
	scheme := r.URL.Scheme

	if len(strings.TrimSpace(scheme)) == 0 {
		scheme = "http"
	}

	if len(strings.TrimSpace(hostname)) == 0 {
		pURL, err := url.Parse(scheme + "://" + r.Host)
		if err != nil {
			log.Fatal("url error: ", err)
			return "", ""
		}
		hostname = pURL.Hostname()
	}
	return scheme, hostname
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//t.Execute(w, p)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		r.ParseForm()       // parse arguments, you have to call this by yourself
		fmt.Println(r.Form) // print form information in server side
		fmt.Println("path", r.URL.Path)
		fmt.Println("scheme", r.URL.Scheme)
		fmt.Println(r.Form["url_long"])
		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
		fn(w, r, m[2])
	}
}

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	GenerateNewCookie(w, "testuid", "aaa")
	GetUserCookie(r, "testuid")
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func TestDBwrite() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)
	defer db.Close()

	fmt.Println("# Inserting values")

	var lastInsertId int
	db.QueryRow("INSERT INTO userinfo(username,created) VALUES($1,$2) returning uid;", "test1", time.Now()).Scan(&lastInsertId)
	fmt.Println("test1 last inserted id =", lastInsertId)

	db1 := GetDBHandle()
	db1.QueryRow("INSERT INTO userinfo(username,created) VALUES($1,$2) returning uid;", "test2", time.Now()).Scan(&lastInsertId)
	fmt.Println("test2 last inserted id =", lastInsertId)
	//db1.QueryRow("INSERT INTO roleassignment(username,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", "aaa", "sysadminrole", time.Now(), "aaa").Scan(&lastInsertId)
	//fmt.Println("roleassign last inserted id =", lastInsertId)

}

func main() {
	//p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page")}
	//p1.save()
	//p2, _ := loadPage("TestPage")
	//fmt.Println(string(p2.Body))

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)
	defer db.Close()

	fmt.Println("# Inserting values")

	var lastInsertId int
	db.QueryRow("INSERT INTO userinfo(username,created) VALUES($1,$2) returning uid;", "aaaaaa1", time.Now()).Scan(&lastInsertId)
	fmt.Println("last inserted id =", lastInsertId)

	TestDBwrite()

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	//http.HandleFunc("/hobby/", hobbyHandler)
	http.HandleFunc("/login", login)

	//http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/www2"))))
	http.Handle("/hello/", helloHandler{})
	//http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www")))) //do not add slash at the beginning of http.Dir path, or add "./" to indicate current folder. Otherwise the path will not be found.
	//http.HandleFunc("/", defaultHandler)

	//InitDBTablewithValue()

	AddUserPermHandler()
	http.ListenAndServe(":"+ServerPort, nil)

}
