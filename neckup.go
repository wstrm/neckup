package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Flag types
var (
	flagTitle       string
	flagPageURI     string
	flagFileURI     string
	flagListenPort  string
	flagUploadDir   string
	flagTmpDir      string
	flagRandPrefix  int
	flagFilenameLen int
)

// init function initializes all the flags for later usage.
// All flags can be user defined, but also they also have a
// default value.
func init() {

	// Constant flags
	const (
		defaultFlagTitle       = "neckup"
		defaultFlagPageURI     = "http://yourdomain.com"
		defaultFlagFileURI     = "http://files.yourdomain.com"
		defaultFlagListenPort  = "8080"
		defaultFlagUploadDir   = "./files"
		defaultFlagRandPrefix  = 24
		defaultFlagFilenameLen = 6
	)

	// Variable flags
	var defaultFlagTmpDir = os.TempDir()

	flag.StringVar(&flagTitle, "title", defaultFlagTitle, "the title that is shown in the view")
	flag.StringVar(&flagPageURI, "page_uri", defaultFlagPageURI, "the page URI that is used in the view")
	flag.StringVar(&flagFileURI, "file_uri", defaultFlagFileURI, "the file URI where the user can find the files")
	flag.StringVar(&flagListenPort, "port", defaultFlagListenPort, "port that the server shoud listen to")
	flag.StringVar(&flagUploadDir, "upload_dir", defaultFlagUploadDir, "directory that the server should save all uploaded files to")
	flag.StringVar(&flagTmpDir, "tmp_dir", defaultFlagTmpDir, "directory that the server should temporarily store file uploads")
	flag.IntVar(&flagRandPrefix, "rand_prefix", defaultFlagRandPrefix, "length of random string that prefixes the temporary filename upon upload")
	flag.IntVar(&flagFilenameLen, "filename_len", defaultFlagFilenameLen, "length of the base filename (excluding extension)")

	flag.Parse()
}

var (
	views      = template.Must(template.ParseFiles("views/index.html"))         // Cache all the templates
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // Allowed characters for random string generator @see randomString
)

// viewHandler renders views/templates.
//
// The function takes three arguments, resWriter which should contain
// the response to write to, view which should contain the view excluding its
// extension (index, upload etc. And not index.html, upload.html etc.).

// Lastly the data argument takes an interface that is optional and can contain
// data that should also be sent to the view.
func viewHandler(resWriter http.ResponseWriter, view string, data interface{}) {

	page := struct {
		Title   string
		PageURI string
		FileURI string
		Data    interface{}
	}{
		flagTitle,
		flagPageURI,
		flagFileURI,
		data,
	}

	views.ExecuteTemplate(resWriter, view+".html", &page)

}

// uploadHandler handles upload requests.
//
// If the request from the client is of the type GET, it'll call
// viewHandler and thereby render the index page.
//
// Else if the request from the client is of the type POST, it'll
// upload all files contained in the request and then call viewHandler
// and thereby render the index with a populated data parameter containing
// the status.
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
			randFilenamePart := randomString(flagRandPrefix)
			fileHash := md5.New()
			part, err := reader.NextPart()

			if err == io.EOF {
				break // Done
			}

			if part.FileName() == "" {
				continue // Empty file name, skip current iteration
			}

			tempPath := filepath.Join(flagTmpDir, randFilenamePart+part.FileName())
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

			finalFilename := hex.EncodeToString(fileHash.Sum(nil))[0:flagFilenameLen] + filepath.Ext(tempPath)
			finalFilepath := filepath.Join(flagUploadDir, finalFilename)

			// Do not copy to storage path if file already exist
			if _, err := os.Stat(finalFilepath); os.IsNotExist(err) {
				os.Rename(tempPath, finalFilepath)
			}

			files[finalFilename] = part.FileName()

		}

		viewHandler(resWriter, "index", files)

	default:
		resWriter.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// randomString generates a random string and returns it.
//
// The length argument decides the length of the random string that should
// be generated.
func randomString(length int) string {

	bits := make([]rune, length)

	for char := range bits {
		bits[char] = characters[rand.Intn(len(characters))]
	}

	return string(bits)
}

// main function initializes everything.
func main() {

	// Seed pseudo-random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", uploadHandler)

	err := http.ListenAndServe(":"+flagListenPort, nil)

	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

}
