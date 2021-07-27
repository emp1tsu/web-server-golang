package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

var tpl *template.Template

type user struct {
	UserName string
	Password string
	First    string
	Last     string
}

var dbSessions = map[string]string{}
var dbUsers = map[string]user{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/vip", vip)
	http.HandleFunc("/signup", signup)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe("localhost:8080", nil)

}

func index(w http.ResponseWriter, req *http.Request) {
	var u user

	c, err := req.Cookie("session")
	if err != nil {
		u = user{}
	}

	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}

	fmt.Printf("%T %v\n", u, u)

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func vip(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "vip.gohtml", nil)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		sID := uuid.New()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		u := user{un, p, f, l}
		dbUsers[un] = u

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}
