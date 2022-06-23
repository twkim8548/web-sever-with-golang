package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadTest(t *testing.T) {
	assert.New(t)

	path := "/Users/taewoong/Desktop/S21.pdf"
	file, _ := os.Open(path)
	defer file.Close()

	os.RemoveAll("./uploads")

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))

	assert.NoError(t, err)

	io.Copy(multi, file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)

	req.Header.Set("Content-type", writer.FormDataContentType())

	uploadsHandler(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, err = os.Stat(uploadFilePath)
	assert.NoError(t, err)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)

	uploadData := []byte{}
	originData := []byte{}

	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(t, originData, uploadData)
}
