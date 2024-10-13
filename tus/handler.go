package tus

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/LinuxSploit/MediaBucket/transcoder"
	"github.com/LinuxSploit/MediaBucket/utils"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
)

var (
	ImageFileTypes map[string]bool = map[string]bool{
		"image/png":  true,
		"image/webp": true,
		"image/jpeg": true,
	}

	VideoFileTypes map[string]bool = map[string]bool{
		"video/mp4":        true,
		"video/webm":       true,
		"video/quicktime":  true,
		"video/avi":        true,
		"video/x-matroska": true,
	}
)

// SetupTusHandler initializes the tusd handler for managing video uploads
func SetupTusVideoHandler(basePath, storageDir string) (*handler.Handler, error) {
	store := filestore.New(storageDir)
	locker := filelocker.New(storageDir)
	composer := handler.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	tusdHandler, err := handler.NewHandler(handler.Config{
		BasePath:              basePath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		NotifyUploadProgress:  true,
		DisableDownload:       true,
		Cors: &handler.CorsConfig{
			Disable:          false,
			AllowOrigin:      regexp.MustCompile(".*"),
			AllowCredentials: false,
			AllowMethods:     "POST, HEAD, PATCH, OPTIONS, GET, DELETE",
			AllowHeaders:     "Authorization, Origin, X-Requested-With, X-Request-ID, X-HTTP-Method-Override, Content-Type, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Incomplete, Upload-Complete, Upload-Draft-Interop-Version",
			MaxAge:           "86400",
			ExposeHeaders:    "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Incomplete, Upload-Complete, Upload-Draft-Interop-Version",
		},
		PreUploadCreateCallback: func(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
			// Extract session token from the headers
			// sessionToken := hook.HTTPRequest.Header.Get("Authorization")

			// // Validate the session token (e.g., check it against your auth service or database)
			// if sessionToken == "" || middleware.ValidateSessionAndPerm(sessionToken) == 0 {
			// 	return handler.HTTPResponse{
			// 		StatusCode: http.StatusUnauthorized,
			// 		Body:       "Invalid or missing session token",
			// 	}, handler.FileInfoChanges{}, nil
			// }

			fileType, ok := hook.Upload.MetaData["filetype"]
			if !ok {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "missing filetype",
				}, handler.FileInfoChanges{}, nil
			}

			isAllowed, ok := VideoFileTypes[fileType]
			if !ok || !isAllowed {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "Invalid filetype",
				}, handler.FileInfoChanges{}, nil
			}

			// you can add additional metadata to the FileInfo if needed
			newMeta := hook.Upload.MetaData
			newMeta["createdDate"] = time.Now().UTC().Format(time.RFC3339) // Add CreatedDate

			fileInfoChanges := handler.FileInfoChanges{
				MetaData: newMeta,
			}

			return handler.HTTPResponse{}, fileInfoChanges, nil
		},
		PreFinishResponseCallback: func(hook handler.HookEvent) (handler.HTTPResponse, error) {

			log.Println("adding upload to queue for transcoding", hook.Upload.IsFinal)
			transcoder.TranscodeQueue <- hook.Upload.ID
			return handler.HTTPResponse{}, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tusd handler: %w", err)
	}

	// Start a goroutine to handle completed uploads
	go func() {
		for {
			event := <-tusdHandler.CompleteUploads
			log.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	return tusdHandler, nil
}

// SetupTusHandler initializes the tusd handler for managing image uploads
func SetupTusImageHandler(basePath, storageDir string) (*handler.Handler, error) {
	store := filestore.New(storageDir)
	locker := filelocker.New(storageDir)
	composer := handler.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	tusdHandler, err := handler.NewHandler(handler.Config{
		BasePath:              basePath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		NotifyUploadProgress:  true,
		Cors: &handler.CorsConfig{
			Disable:          false,
			AllowOrigin:      regexp.MustCompile(".*"),
			AllowCredentials: false,
			AllowMethods:     "POST, HEAD, PATCH, OPTIONS, GET, DELETE",
			AllowHeaders:     "Authorization, Origin, X-Requested-With, X-Request-ID, X-HTTP-Method-Override, Content-Type, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Incomplete, Upload-Complete, Upload-Draft-Interop-Version",
			MaxAge:           "86400",
			ExposeHeaders:    "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Incomplete, Upload-Complete, Upload-Draft-Interop-Version",
		},
		PreUploadCreateCallback: func(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
			// // Extract session token from the headers
			// sessionToken := hook.HTTPRequest.Header.Get("Authorization")

			// // Validate the session token (e.g., check it against your auth service or database)
			// if sessionToken == "" || middleware.ValidateSessionAndPerm(sessionToken) == 0 {
			// 	return handler.HTTPResponse{
			// 		StatusCode: http.StatusUnauthorized,
			// 		Body:       "Invalid or missing session token",
			// 	}, handler.FileInfoChanges{}, nil
			// }

			fileType, ok := hook.Upload.MetaData["filetype"]
			if !ok {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "missing filetype",
				}, handler.FileInfoChanges{}, nil
			}

			isAllowed, ok := ImageFileTypes[fileType]
			if !ok || !isAllowed {
				fmt.Println(fileType)
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "Invalid filetype",
				}, handler.FileInfoChanges{}, nil
			}

			// you can add additional metadata to the FileInfo if needed
			newMeta := hook.Upload.MetaData
			newMeta["createdDate"] = time.Now().UTC().Format(time.RFC3339) // Add CreatedDate

			fileInfoChanges := handler.FileInfoChanges{
				MetaData: newMeta,
			}

			return handler.HTTPResponse{}, fileInfoChanges, nil
		},
		PreFinishResponseCallback: func(hook handler.HookEvent) (handler.HTTPResponse, error) {

			log.Println("generate thumbnail of 500 width: ", "/storage/MediaBucket/images/"+hook.Upload.ID, "/storage/MediaBucket/thumbnail/"+hook.Upload.ID+"-500w.webp")
			err := utils.ResizeAndConvertToWebP("/storage/MediaBucket/images/"+hook.Upload.ID, "/storage/MediaBucket/thumbnail/"+hook.Upload.ID+"-500w.webp", 500)
			if err != nil {
				fmt.Println("thumbnail generation err: ", err)
			}
			return handler.HTTPResponse{}, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tusd handler: %w", err)
	}

	// Start a goroutine to handle completed uploads
	go func() {
		for {
			event := <-tusdHandler.CompleteUploads
			log.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	return tusdHandler, nil
}
