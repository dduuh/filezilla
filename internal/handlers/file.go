package handlers

import (
	"context"
	"filezilla/internal/domain"
	minioclient "filezilla/pkg/minclient"
	"filezilla/pkg/responses"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	tenMB      = 10 << 20
	bucketName = "uploads"
)

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.HTTPError(w, "incorrect HTTP method", http.StatusBadRequest)
		return
	}

	// Getting User ID
	userId, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		responses.HTTPError(w, "ID not found", http.StatusUnauthorized)
		return
	}

	// Parsing Form
	if err := r.ParseMultipartForm(tenMB); err != nil { // 10 MB
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Form File
	file, handler, err := r.FormFile("file")
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// File Size
	fileSize := handler.Size

	// Checking if exists
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		logrus.Fatal("missing file: 'handler'")
	}

	// Timestamp
	timestamp := time.Now().Format("20060102150405")
	// Filename
	fileName := fmt.Sprintf("%d_%s%s", userId, timestamp, ext)

	ctx := context.Background()

	// Minio storage
	mio, err := minioclient.NewMinioClient()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// Creating bucket if not exists
	mio.NewBucket(ctx, bucketName)

	// Get Content-Type
	contentType := handler.Header.Get("Content-Type")

	// Uploading object/file
	if err := mio.UploadFile(ctx, bucketName, fileName, file, fileSize, contentType); err != nil {
		logrus.Fatal(err.Error())
	}

	// Final storage url to store to PostgreSQL
	storageUrl := mio.GenerateStorageUrl(bucketName, fileName)

	// Finish File Data
	fileData := domain.File{
		UserId:     userId,
		FileSize:   fileSize,
		CreatedAt:  time.Now(),
		StorageUrl: storageUrl,
	}

	// Uploading file info to PostgreSQL
	fileId, err := h.services.Files.CreateFile(fileData)
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.HTTPResponse(w, http.StatusCreated, "file_id", fileId)
}

func (h *Handler) getFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.HTTPError(w, "incorrect HTTP method", http.StatusBadRequest)
		return
	}

	files, err := h.services.Files.GetFiles()
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.HTTPResponse(w, http.StatusOK, "files", files)
}
