package main

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Settings
const TITLE = "u.wiol.io"           // Title for views
const PAGE_URI = "http://u.wiol.io" // URI for the home page
const FILE_URI = "http://f.wiol.io" // URI for the files
const LISTEN_PORT = "8080"          // The port the server should listen to
const UPLOAD_DIR = "./files/"       // Save all files to this directory

// Length of random string that prefixes the filename upon upload
const TMP_FILENAME_LEN = 24

// Length of the base filename (excluding extension)
const FINAL_FILENAME_LEN = 6

// Cache all the templates
var views = template.Must(template.ParseFiles("views/index.html"))

// Allowed characters for random string generator @see randomString
var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

/**
 * viewHandler
 * Render views/templates.
 *
 * @param		resWriter	The response writer from net/http lib.
 * @param		view			Name for the view that should be rendered.
 * @param		data			Optional, data that can be used by the view.
 * @return						nil
 */
func viewHandler(resWriter http.ResponseWriter, view string, data interface{}) {

	page := struct {
		Title   string
		PageURI string
		FileURI string
		Data    interface{}
	}{
		TITLE,
		PAGE_URI,
		FILE_URI,
		data,
	}

	views.ExecuteTemplate(resWriter, view+".html", &page)

}

/**
 * uploadHandler
 * Handle upload requests.
 *
 * If the request from the client is of the type GET, it'll call
 * viewHandler and thereby render the index page.
 *
 * Else if the request from the client is of the type POST, it'll
 * upload all files contained in the request and then call viewHandler
 * and thereby render the index with a populated data parameter containing
 * the status.
 *
 * @param		resWriter	The response writer from net/http lib.
 * @param		req				The request that'll uploadHandler will respond for.
 * @return						nil
 */
func uploadHandler(resWriter http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "GET":
		viewHandler(resWriter, "index", nil)

	case "POST":
		reader, err := req.MultipartReader()
		files := make(map[string]string)

		if err != nil {
			log.Print(err)

			http.Error(resWriter, "Failed to read multipart stream.", http.StatusInternalServerError)
			return
		}

		for {
			randFilenamePart := randomString(TMP_FILENAME_LEN)
			fileHash := md5.New()
			part, err := reader.NextPart()

			if err == io.EOF {
				break // Done
			}

			if part.FileName() == "" {
				continue // Empty file name, skip current iteration
			}

			tempPath := filepath.Join(os.TempDir(), randFilenamePart+part.FileName())
			tempDest, err := os.Create(tempPath)
			defer tempDest.Close()

			parsedPart := io.TeeReader(part, fileHash) // Feed hash with part

			if err != nil {
				log.Print(err)

				http.Error(resWriter, "Something went wrong.", http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(tempDest, parsedPart); err != nil {
				log.Print(err)

				http.Error(resWriter, "Unable to parse file.", http.StatusInternalServerError)
				return
			}

			finalFilename := hex.EncodeToString(fileHash.Sum(nil))[0:FINAL_FILENAME_LEN] + filepath.Ext(tempPath)
			os.Rename(tempPath, filepath.Join(UPLOAD_DIR, finalFilename))
			files[finalFilename] = part.FileName()

		}

		viewHandler(resWriter, "index", files)

	default:
		resWriter.WriteHeader(http.StatusMethodNotAllowed)
	}

}

/**
 * randomString
 * Generate random string.
 *
 * @param		length	Length of the random that should be generated.
 * @return					Random string with the length specified.
 */
func randomString(length int) string {

	randBits := make([]rune, length)
	for char := range randBits {
		randBits[char] = characters[rand.Intn(len(characters))]
	}

	return string(randBits)
}

/**
 * main
 * Da glorious main function that initializes everything.
 *
 * @return nil
 */
func main() {

	// Seed pseudo-random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", uploadHandler)

	err := http.ListenAndServe(":"+LISTEN_PORT, nil)

	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

}
