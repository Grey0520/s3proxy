package handlers

import (
	"net/http"
	"strings"

	s "github.com/Grey0520/s3proxy/internal/server"
	"github.com/Grey0520/s3proxy/internal/storage"
	"github.com/labstack/echo/v4"
)

type ObjectHandlers struct {
	server *s.Server
}

func NewObjectHandlers(server *s.Server) *ObjectHandlers {
	return &ObjectHandlers{server: server}
}

func (h *ObjectHandlers) GetObject(c echo.Context) error {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")

	storage := *h.server.Storage
	obj, err := storage.GetObject(bucketName, objectName)
	if err != nil {
		return c.XML(http.StatusInternalServerError, err.Error())
	}

	// c.Response().Header().Set("Content-Type", obj.ContentType)
	// c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", obj.Size))
	// c.Response().Header().Set("Last-Modified", obj.LastModified.UTC().Format(http.TimeFormat))

	return c.Stream(http.StatusOK, obj.ContentType, obj.Data)
}

func (h *ObjectHandlers) PutObject(c echo.Context) error {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")

	stg := *h.server.Storage

	// 处理 Copy Object 请求
	if srcPath := c.Request().Header.Get("x-amz-copy-source"); len(srcPath) != 0 {
		srcPath := c.Request().Header.Get("x-amz-copy-source")
		trimmedPath := strings.TrimPrefix(srcPath, "/")
		parts := strings.Split(trimmedPath, "/")

		if len(parts) != 2 {
			return c.XML(http.StatusBadRequest, "Invalid copy source")
		}

		desBucketName := c.Param("bucketName")
		desObjectName := c.Param("objectName")
		srcBucketName := parts[0]
		srcObjectName := parts[1]

		stg := *h.server.Storage
		err := stg.CopyObject(srcBucketName, srcObjectName, desBucketName, desObjectName)
		if err != nil {
			return c.XML(http.StatusInternalServerError, err.Error())
		}
		return nil
	}

	// 剩下的是从请求体中读取数据的请求
	obj := &storage.Object{
		ContentType: c.Request().Header.Get("Content-Type"),
		Data:        c.Request().Body,
	}
	err := stg.PutObject(bucketName, objectName, obj)
	if err != nil {
		return c.XML(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *ObjectHandlers) DeleteObject(c echo.Context) error {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")

	stg := *h.server.Storage
	err := stg.DeleteObject(bucketName, objectName)
	if err != nil {
		return c.XML(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *ObjectHandlers) CopyObject(c echo.Context) error {
	// 处理类似 “/bucketName/objectName” 的路径
	srcPath := c.Request().Header.Get("x-amz-copy-source")
	trimmedPath := strings.TrimPrefix(srcPath, "/")
	parts := strings.Split(trimmedPath, "/")

	if len(parts) != 2 {
		return c.XML(http.StatusBadRequest, "Invalid copy source")
	}

	desBucketName := c.Param("bucketName")
	desObjectName := c.Param("objectName")
	srcBucketName := parts[0]
	srcObjectName := parts[1]

	stg := *h.server.Storage
	err := stg.CopyObject(srcBucketName, srcObjectName, desBucketName, desObjectName)
	if err != nil {
		return c.XML(http.StatusInternalServerError, err.Error())
	}
	return nil
}
