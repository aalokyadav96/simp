package main

import (
    "fmt"
    "net/http"
    "log"
	"io"
	"os"
	"time"
	"os/exec"
	"crypto/rand"
	rndm "math/rand"
	"crypto/md5"
	"path/filepath"

    "github.com/julienschmidt/httprouter"
)


func HasAuthCookie(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLoggedIn(w,r) {
			next(w, r, ps)
		} else {
			fmt.Println("nah")
			fmt.Fprintf(w,"nah")
		}
	}
}

func ViewUser() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
    }
}

func Me() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprintf(w, "meme")
    }
}


func UploadFileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isLoggedIn(w,r) {
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
	} else if r.Method == "GET" && isLoggedIn(w,r) {
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "upload.html", nil)
	} else if r.Method == "POST" && isLoggedIn(w,r) {
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
		}
		var fileEndings string
		var fileName string

		files := r.MultipartForm.File["imgfile"]
		for _, fileHeader := range files {
			log.Println("hao")
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			defer file.Close()
			log.Println("file OK")
			// Get and print out file size
			fileSize := fileHeader.Size
			fmt.Printf("File size (bytes): %v\n", fileSize)
			// validate file size
			if fileSize > maxUploadSize {
				renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			}
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				renderError(w, "INVALID_FILE"+http.DetectContentType(fileBytes), http.StatusBadRequest)
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			detectedFileType := http.DetectContentType(fileBytes)
			switch detectedFileType {
			case "video/mp4":
				fileEndings = ".mp4"
				break
			case "video/webm":
				fileEndings = ".webm"
				break
			case "image/gif":
				fileEndings = ".gif"
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			}
			fileName = GenerateName(randToken(12))
			//		fileEndings, err := mime.ExtensionsByType(detectedFileType)

			if err != nil {
				renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			}
			newFileName := fileName + fileEndings

			newPath := filepath.Join(uploadPath, newFileName)
			fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
			//			ffmpeg -i %1 -filter:v scale=-2:640:flags=lanczos -c:a copy -pix_fmt yuv420p %1_lcz.mp4
			go FFConvert(fileName , fileEndings )
			http.Redirect(w, r, "/v/"+fileName, http.StatusSeeOther)
			}
			
//				tmpl.ExecuteTemplate(w, "show.html", fileName)
	}
}

func FFConvert(fileName string, fileEndings string) {
	getFrom := uploadPath + "/" + fileName + fileEndings
	saveAs := streamPath + "/" + fileName + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", getFrom, "-filter:v", "scale=-2:640:flags=lanczos", "-c:a", "copy", "-pix_fmt", "yuv420p", saveAs)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

//~ func FFConvert(fileName string, fileEndings string)  {
//~ getFrom := uploadPath+"/"+fileName+fileEndings
//~ saveAs := uploadPath+"/"+fileName+"_lcz.mp4"
//~ cmd := exec.Command("ffmpeg", "-i", getFrom, "-filter:v", "scale=-2:640:flags=lanczos", "-c:a", "copy", "-pix_fmt", "yuv420p", saveAs)
//~ err := cmd.Start()
//~ if err != nil {
//~ log.Fatal(err)
//~ }
//~ log.Printf("Waiting for command to finish...")
//~ err = cmd.Wait()
//~ log.Printf("Command finished with error: %v", err)
//~ os.Remove(getFrom)
//~ }


func ViewPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println(r.URL.Path)
		tmpl.ExecuteTemplate(w, "show.html", ps.ByName("PostId"))
	}
}

func ViewAllFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "viewall.html", SearchFiles(uploadPath))
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(r.URL.Path)
	if r.Method == "DELETE" {
		fileName := ps.ByName("PostId")
		fmt.Println(uploadPath + "/" + fileName)
		os.Remove(uploadPath + "/" + fileName)
		fmt.Println(streamPath + "/" + fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".mp4")
		os.Remove(streamPath + "/" + fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".mp4")
		XHRrespond(w, "Deleted")
	}
}

func SearchFiles(dir string) []string {
	var allFiles []string
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		switch file.Name() {
		case "$RECYCLE.BIN", "$Recycle.Bin":
			break
		case "System Volume Information":
			break
		default:
			allFiles = append(allFiles, file.Name())
		}
	}
	return allFiles
}

func GenerateName(w int64) string {
	rndm.Seed(time.Now().Unix()) // initialize global pseudo random generator
	p1 := fmt.Sprintf(adjectives[rndm.Intn(len(adjectives))])
	p2 := fmt.Sprintf(adjectives[rndm.Intn(len(adjectives))])
	p3 := fmt.Sprintf(animals[rndm.Intn(len(animals))])
	return p1 + p2 + p3
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func randToken(len int) int64 {
	b := make([]byte, len)
	n, _ := rand.Read(b)
	return int64(n)
}

func XHRrespond(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, message)
}

func HashIt(strToHash string) string {
	data := []byte(strToHash)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func Ignore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "favicon.ico")
}