package main

import (
    "net/http"
    "log"
	"html/template"

    "github.com/julienschmidt/httprouter"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

const maxUploadSize = 10 * 1024 * 1024 // 10 mb
const uploadPath = "./uploads"
const streamPath = "./cmp"
var postersDir = "./posters"
var userpicPath = "./userpic"

func main() {
    router := httprouter.New()
    router.GET("/", HasAuthCookie(UploadFileHandler))

	router.GET("/search", Search)
	router.GET("/tag/:tag", Tags)

	router.GET("/login", loginHandler)
    router.POST("/login", loginHandler)
    router.GET("/register", registerHandler)
    router.POST("/register", registerHandler)
    router.POST("/logout", HasAuthCookie(logoutHandler))

	//~ router.GET("/hello/:name", Hello)
    router.GET("/@:name", HasAuthCookie(ViewUser()))
    router.GET("/me", HasAuthCookie(Me()))
    router.GET("/manage", HasAuthCookie(ManageContent()))
	
	router.GET("/upload", HasAuthCookie(UploadFileHandler))
	router.POST("/upload", HasAuthCookie(UploadFileHandler))
	router.GET("/v/:PostId", ViewPost)
	router.GET("/viewall", HasAuthCookie(ViewAllFiles))
	router.DELETE("/del/:PostId", HasAuthCookie(DeleteFile))
	
	router.GET("/fav/favicon.ico", Ignore)
	
	static := httprouter.New()
	static.ServeFiles("/files/*filepath", http.Dir(streamPath))
	static.ServeFiles("/giant/*filepath", http.Dir(uploadPath))
	static.ServeFiles("/assets/*filepath", http.Dir("./assets"))
	static.ServeFiles("/poster/*filepath", http.Dir(postersDir))
	static.ServeFiles("/userpic/*filepath", http.Dir(userpicPath))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
//	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.NotFound = static


	log.Println("Starting Server")
    log.Fatal(http.ListenAndServe("localhost:4000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tmpl.ExecuteTemplate(w,"head.html",nil)
	tmpl.ExecuteTemplate(w,"index.html",nil)
}
