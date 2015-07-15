package main

import (
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Flag types
var (
	flagTitle         string
	flagPageURI       string
	flagFileURI       string
	flagListenPort    string
	flagUploadDir     string
	flagTmpDir        string
	flagIndexView     string
	flagDisallowChars string
	flagRandPrefix    int
	flagFilenameLen   int
	flagVersion       bool
)

// init function initializes all the flags for later usage.
// All flags can be user defined, but also they also have a
// default value.
func init() {

	// Bump this upon version change (semver)
	const currentVersion = "0.0.3"

	// Constant flags
	const (
		defaultFlagTitle         = "neckup"
		defaultFlagPageURI       = "http://yourdomain.com"
		defaultFlagFileURI       = "http://files.yourdomain.com"
		defaultFlagListenPort    = "8080"
		defaultFlagUploadDir     = "./files"
		defaultFlagIndexView     = "minimal"
		defaultFlagDisallowChars = "lIO0-"
		defaultFlagRandPrefix    = 24
		defaultFlagFilenameLen   = 6
		defaultFlagVersion       = false
	)

	// Variable flags
	defaultFlagTmpDir := os.TempDir()

	flag.StringVar(&flagTitle, "title", defaultFlagTitle, "the title that is shown in the view")
	flag.StringVar(&flagPageURI, "page_uri", defaultFlagPageURI, "the page URI that is used in the view")
	flag.StringVar(&flagFileURI, "file_uri", defaultFlagFileURI, "the file URI where the user can find the files")
	flag.StringVar(&flagListenPort, "port", defaultFlagListenPort, "port that the server shoud listen to")
	flag.StringVar(&flagUploadDir, "upload_dir", defaultFlagUploadDir, "directory that the server should save all uploaded files to")
	flag.StringVar(&flagTmpDir, "tmp_dir", defaultFlagTmpDir, "directory that the server should temporarily store file uploads")
	flag.StringVar(&flagIndexView, "index_view", defaultFlagIndexView, "index view to show on root page")
	flag.StringVar(&flagDisallowChars, "disallow_chars", defaultFlagDisallowChars, "disallowed characters for final filenames")
	flag.IntVar(&flagRandPrefix, "rand_prefix", defaultFlagRandPrefix, "length of random string that prefixes the temporary filename upon upload")
	flag.IntVar(&flagFilenameLen, "filename_len", defaultFlagFilenameLen, "length of the base filename (excluding extension)")
	flag.BoolVar(&flagVersion, "v", defaultFlagVersion, "print current neckup version and exit")

	flag.Parse()

	// If version flag is true, print version and exit
	if flagVersion {
		fmt.Println(currentVersion)
		os.Exit(0)
	}

	return
}

var (
	// Cache (all) the template(s)
	views = template.Must(template.ParseGlob(filepath.Join("./views/", "*.html")))

	// Allowed characters for random string generator @see randomString
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// viewHandler renders views/templates.
//
// The function takes three arguments, writer which should contain
// the response to write to, view which should contain the view excluding its
// extension (index, upload etc. And not index.html, upload.html etc.).

// Lastly the data argument takes an interface that is optional and can contain
// data that should also be sent to the view.
func viewHandler(writer http.ResponseWriter, view string, data interface{}) {

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

	err := views.ExecuteTemplate(writer, view+".html", &page)

	if err != nil {
		log.Print(err)

		http.Error(writer, "Failed to compile view.", http.StatusInternalServerError)
		return
	}

	return
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
func uploadHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		switch request.Method {

		case "POST":

			reader, err := request.MultipartReader()
			files := make(map[string]string)

			if err != nil {
				log.Print(err)

				http.Error(writer, "Failed to read multipart stream.", http.StatusInternalServerError)
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

					http.Error(writer, "Something went wrong.", http.StatusInternalServerError)
					return
				}

				if _, err := io.Copy(tempDest, parsedPart); err != nil {
					log.Print(err)

					http.Error(writer, "Unable to parse file.", http.StatusInternalServerError)
					return
				}

				finalFilename := stripChars(base64.URLEncoding.EncodeToString(fileHash.Sum(nil)), flagDisallowChars)[0:flagFilenameLen] + filepath.Ext(tempPath)
				finalFilepath := filepath.Join(flagUploadDir, finalFilename)

				// Do not copy to storage path if file already exist
				if _, err := os.Stat(finalFilepath); os.IsNotExist(err) {
					os.Rename(tempPath, finalFilepath)
				} else { // Remove temporary file
					os.Remove(tempPath)
				}

				files[finalFilename] = part.FileName()

			}

			viewHandler(writer, flagIndexView, files)

		case "GET":
			viewHandler(writer, flagIndexView, nil)

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
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

// stripChars strips unwanted characters from a string.
//
// The return value is either the stripped string or -1 if
// no chars was declared.
func stripChars(str, chars string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chars, r) < 0 {
			return r
		}

		return -1
	}, str)
}

// main function initializes everything.
func main() {

	// Seed pseudo-random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	http.Handle("/", uploadHandler())

	err := http.ListenAndServe(":"+flagListenPort, nil)

	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	return
}
