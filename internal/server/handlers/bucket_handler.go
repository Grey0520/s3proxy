package handlers

import (
	"net/http"

	s "github.com/Grey0520/s3proxy/internal/server"
	"github.com/labstack/echo/v4"
)

type BucketHandler struct {
	server *s.Server
}

func NewBucketHandlers(server *s.Server) *BucketHandler {
	return &BucketHandler{server: server}
}

func (h *BucketHandler) CreateBucket(c echo.Context) error {
	bucketName := c.Param("bucketName")

	stg := *h.server.Storage
	err := stg.CreateBucket(bucketName)
	if err != nil {
		return c.XML(http.StatusInternalServerError, err.Error())
	}

	return c.XML(http.StatusOK, "Bucket created")
}
