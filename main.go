package main

import (
    "net/http"
    "log"
	"html/template"

    "github.com/julienschmidt/httprouter"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

const maxUploadSize = 10 * 1024 * 1024 // 8 mb
const uploadPath = "./uploads"
const streamPath = "./cmp"

var userpicPath = "./userpic"

func main() {
    router := httprouter.New()
    router.GET("/", Index)

	router.GET("/search", Index)

	router.GET("/login", loginHandler)
    router.POST("/login", loginHandler)
    router.GET("/register", registerHandler)
    router.POST("/register", registerHandler)
    router.POST("/logout", logoutHandler)

	//~ router.GET("/hello/:name", Hello)
    router.GET("/@:name", HasAuthCookie(ViewUser()))
    router.GET("/me", HasAuthCookie(Me()))
	
	router.GET("/upload", UploadFileHandler)
	router.POST("/upload", UploadFileHandler)
	router.GET("/v/:PostId", ViewPost)
	router.GET("/viewall", ViewAllFiles)
	router.DELETE("/del/:PostId", DeleteFile)
	
	router.GET("/favicon.ico", Ignore)
	
	static := httprouter.New()
	static.ServeFiles("/files/*filepath", http.Dir(streamPath))
	static.ServeFiles("/giant/*filepath", http.Dir(uploadPath))
	static.ServeFiles("/userpic/*filepath", http.Dir(userpicPath))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.NotFound = static


	log.Println("Starting Server")
    log.Fatal(http.ListenAndServe("localhost:4000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tmpl.ExecuteTemplate(w,"head.html",nil)
	tmpl.ExecuteTemplate(w,"index.html",nil)
}

//~ func Hello() httprouter.Handle {
    //~ return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        //~ fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
    //~ }
//~ }