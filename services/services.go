package services

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MaxUploadSize -
const MaxUploadSize = 100 * 1024 * 1024

// WasmDir -
var WasmDir = "output/"

// UploadPath -
var UploadPath = "test/"

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

// UploadFile -
func UploadFile(c *gin.Context) {
	// Maximum upload of 10 MB files
	w := c.Writer
	suffix := fmt.Sprintf("%v", time.Now().Unix())
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, MaxUploadSize)
	if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		// fmt.Println("Error Retrieving the File")
		// 	fmt.Println(err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// fileBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	renderError(w, "INVALID_FILE", http.StatusBadRequest)
	// 	return
	// }
	// Get handler for filename, size and headers
	// file, handler, err := r.FormFile("myFile")
	// if err != nil {
	// 	fmt.Println("Error Retrieving the File")
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	_, err = os.Stat(UploadPath + suffix)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(UploadPath+suffix, 0755)
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
	dst, err := os.Create(UploadPath + handler.Filename)
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
	wasmFileName := strings.Split(handler.Filename, ".")[0] + "_" + suffix + ".wasm"

	cmd1 := exec.Command("go", "mod", "init")
	cmd1.Dir = UploadPath + suffix
	cmd2 := exec.Command("go", "mod", "download")
	cmd2.Dir = UploadPath + suffix
	// cmd2 := exec.Command("./setenv.sh")
	cmd3 := exec.Command("go", "build", "-o", WasmDir+wasmFileName, UploadPath+handler.Filename)
	// cmd3 := exec.Command("go", "build", "-o", "output/"+wasmFileName, UploadPath+handler.Filename)
	cmd3.Env = os.Environ()
	cmd3.Env = append(cmd3.Env, "GOOS=js")
	cmd3.Env = append(cmd3.Env, "GOARCH=wasm")
	err = cmd1.Run()
	if err != nil {
		fmt.Println("1:" + fmt.Sprint(err))
		return
	}
	_, err = cmd2.Output()
	if err != nil {
		fmt.Println("2 " + fmt.Sprint(err))
		return
	}
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
		c.JSON(404, gin.H{
			"status": fmt.Sprint(err),
			"reason": stderr.String(),
		})
		// fmt.Fprintf(w, fmt.Sprint(err)+": "+stderr.String())
		return
	}
	cmd4 := exec.Command("rm", "-rf", UploadPath+suffix)
	_, err = cmd4.Output()
	if err != nil {
		fmt.Println("Error while deleting the test folder")
	}
	outputFile := "getwasm/" + wasmFileName
	c.JSON(200, gin.H{
		"filePath": outputFile,
	})
	// fmt.Println(strconv.Quote(wasmFileName))
	// c.File(outputFile)
	// w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(wasmFileName))
	// w.Header().Set("Content-Type", "application/wasm")
	// http.ServeFile(w, r, outputFile)
	// fmt.Println(string(stdout))
}

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func Display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

// DownloadWasmFile -
func DownloadWasmFile(c *gin.Context) {
	// dir := c.Param("uploadproxy")
	fileID := c.Param("filename")
	// filePath := fmt.Sprintf("%v", c.Request.URL)
	// filePath = WasmDir +strings.Replace(filePath, models.ProjectCFG.ProjectID+"/images/", "/", 1)

	// if isFound {
	filePath := filepath.Join(WasmDir, fileID)
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"status": "Not Found",
			"reason": "File does not exists",
		})
		return
	}
	c.File(filePath)
	// }
}
