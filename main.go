package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/LinuxSploit/MediaBucket/middleware"
	"github.com/LinuxSploit/MediaBucket/transcoder"
	"github.com/LinuxSploit/MediaBucket/tus"
	"github.com/joho/godotenv"
)

func init() {

	// load .env file
	godotenv.Load()

	log.Println("init completed!")
}

func main() {
	log.Println("service started...")
	// Create TUS video Upload handler, init basePath, storageDir
	videoHandler, err := tus.SetupTusVideoHandler(os.Getenv("ServerURL")+os.Getenv("VideoUploadPath"), "/storage/MediaBucket/videos")
	if err != nil {
		log.Fatalf("Unable to create photo handler: %v", err)
	}

	// Create TUS video Upload handler, init basePath, storageDir
	imageHandler, err := tus.SetupTusImageHandler(os.Getenv("ServerURL")+os.Getenv("ImageUploadPath"), "/storage/MediaBucket/images")
	if err != nil {
		log.Fatalf("Unable to create photo handler: %v", err)
	}

	mux := http.NewServeMux()

	// Register TUS video Upload handler to /upload/ route
	mux.Handle("/video/", http.StripPrefix("/video/", videoHandler))

	// Register TUS image Upload handler to /image/ route
	mux.Handle("/image/", http.StripPrefix("/image/", imageHandler))

	// Serve HLS video streams
	videoFileServer := http.StripPrefix("/hls/", NoDirListingFileServer(http.Dir("/storage/MediaBucket/hls/")))
	mux.Handle("/hls/", middleware.CORSMiddleware(videoFileServer))

	//thumbnail server
	thumbnailFileServer := http.StripPrefix("/thumbnail/", NoDirListingFileServer(http.Dir("/storage/MediaBucket/thumbnail/")))
	mux.Handle("/thumbnail/", middleware.CORSMiddleware(thumbnailFileServer))

	// Serve the video upload demo page
	mux.HandleFunc("/video-demo", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/video-demo.html"))
		tmpl.Execute(w, os.Getenv("ServerURL"))
	})

	// Serve the image upload demo page
	mux.HandleFunc("/image-demo", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/image-demo.html"))
		tmpl.Execute(w, os.Getenv("ServerURL"))
	})

	// Start the video transcoding worker
	transcoder.StartTranscodeWorker("/storage/MediaBucket/videos/", "/storage/MediaBucket/hls/")

	// Start the HTTP server
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

// NoDirListingFileServer wraps the http.FileServer to disable directory listings
func NoDirListingFileServer(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the absolute path to prevent directory traversal
		upath := path.Clean(r.URL.Path)

		f, err := root.Open(upath)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		defer f.Close()

		// Get file info
		info, err := f.Stat()
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// deny access to directory path
		if info.IsDir() {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Serve the file
		fs.ServeHTTP(w, r)
	})
}
