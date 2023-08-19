package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "exampleCookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func registerHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	switch request.Method {
	case "GET" : 
		tmpl.ExecuteTemplate(response,"register.html",nil)
	case "POST" : 
	name := request.FormValue("username")
	pass := request.FormValue("password")
	fmt.Println(name, "  " ,pass)
	if name != "" && pass != "" {
		rdxHset("user",name,pass)
		val, _ := rdxHget("user",name)
		fmt.Println(val)
		http.Redirect(response, request, "/login", 302)
	} else {
		http.Redirect(response, request, "/register", 302)
	}
	default:
		fmt.Fprintf(response,"Method Not Allowed")
	}
}

// login handler

// index page
const already = `<!DOCTYPE HTML><p>Already Logged In</p>`

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch r.Method {
	case "GET" : {
		if isLoggedIn(w,r) {
			fmt.Fprintf(w, already)
		} else {
			tmpl.ExecuteTemplate(w,"login.html",nil)
		}
	}
	case "POST" : {
	name := r.FormValue("username")
	pass := r.FormValue("password")
	if name != "" && pass != "" {
		val, _ := rdxHget("user",name)
		fmt.Println(val)
		fmt.Println(pass)
		if val == pass {
			setSession(name,w)
			http.Redirect(w, r, "/", 301)
		} else {
			fmt.Fprintf(w,"Wrong credentials")
		}
	}
	}
	default:
		fmt.Fprintf(w,"Method Not Allowed")
	}
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		cookie := &http.Cookie{
		Name:   "exampleCookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
	http.Redirect(response, request, "/", 302)
}


func indexPageHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	tmpl.ExecuteTemplate(response,"login.html",nil)
}

func setSession(name string, w http.ResponseWriter) {
    cookie := http.Cookie{
        Name:     "exampleCookie",
        Value:    name,
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, &cookie)
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
    cookie, err := r.Cookie("exampleCookie")
    if err != nil {
        return false
    }
	fmt.Println(cookie.Value)
	return true
}