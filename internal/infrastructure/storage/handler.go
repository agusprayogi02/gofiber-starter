package storage

import (
	"mime/multipart"
	"os"
	"starter-gofiber/pkg/apierror"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// dirPath: harus diawali dan diakhiri dengan "/"
func UploadFile(c *fiber.Ctx, file *multipart.FileHeader, dirPath string) (string, error) {
	uuid := uuid.New().String()
	arrFile := strings.Split(file.Filename, ".")
	path := "./public" + dirPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Folder tidak ada, buat folder
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Info("Error creating folder:", err)
			panic(err)
		}
		log.Info("Folder created successfully.")
	}
	fileName := uuid + "." + arrFile[len(arrFile)-1]
	path += fileName
	if err := c.SaveFile(file, path); err != nil {
		return fileName, &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "U-UploadFile",
		}
	}
	return fileName, nil
}

func DeleteFile(fileName *string, dirPath string) error {
	if fileName == nil {
		return nil
	}
	path := "./public" + dirPath + *fileName
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "U-DeleteFile",
		}
	}
	return nil
}
