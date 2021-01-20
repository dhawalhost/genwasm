package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	_, err = os.Stat("test")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("test", 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	_, err = os.Stat("output")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("output", 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	// Create file
	dst, err := os.Create("test/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
	suffix := fmt.Sprintf("%v", time.Now().Unix())
	wasmFileName := strings.Split(handler.Filename, ".")[0] + "_" + suffix + ".wasm"

	// cmd1 := exec.Command("export", "GOOS=js")

	// cmd2 := exec.Command("./setenv.sh")
	fmt.Println("go", "build", "-o", "output/"+wasmFileName, "test/"+handler.Filename)
	cmd3 := exec.Command("go", "build", "-o", "output/"+wasmFileName, "test/"+handler.Filename)
	// cmd3 := exec.Command("go", "build", "-o", "output/"+wasmFileName, "test/"+handler.Filename)
	cmd3.Env = os.Environ()
	cmd3.Env = append(cmd3.Env, "GOOS=js")
	cmd3.Env = append(cmd3.Env, "GOARCH=wasm")
	// stdout, err := cmd1.Output()

	// if err != nil {
	// 	fmt.Println("1 " + err.Error())
	// 	return
	// }
	// _, err = cmd2.Output()

	// if err != nil {
	// 	fmt.Println("2 ")
	// 	return
	// }
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd3.Stdout = &out
	cmd3.Stderr = &stderr
	err = cmd3.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	outputFile := "/home/legerrac/go/src/github.com/dhawalhost/genwasm/output/" + wasmFileName
	fmt.Println(strconv.Quote(wasmFileName))
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(wasmFileName))
	w.Header().Set("Content-Type", "application/wasm")
	http.ServeFile(w, r, outputFile)
	fmt.Println("Result: " + out.String())
	// fmt.Println(string(stdout))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
