package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const maxMemory = int64(1 * 1024 * 1024)

func main() {
	router := gin.Default()
	allowedExt := []string{".jpg", ".jpeg", ".png"}
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = maxMemory // 2 MiB
	router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["files"]

		for _, file := range files {
			isAllowed := false
			for _, ext := range allowedExt {
				if filepath.Ext(file.Filename) == ext {
					isAllowed = true
				}
			}

			if !isAllowed {
				c.JSON(http.StatusBadRequest, gin.H{"error": "file extension not allowed"})
				return
			}
			if file.Size > maxMemory {
				c.JSON(http.StatusBadRequest, gin.H{"error": "file size too big"})
				return
			}
			filename := fmt.Sprintf("%s.%s", "files/", file.Filename)
			//filepath.Base(file.Filename)
			fmt.Println(filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email)
	})
	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}
