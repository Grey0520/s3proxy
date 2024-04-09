package routes

import (
	s "github.com/Grey0520/s3proxy/internal/server"
	"github.com/Grey0520/s3proxy/internal/server/handlers"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureRoutes(server *s.Server) {
	objectHanlder := handlers.NewObjectHandlers(server)
	bucketHandler := handlers.NewBucketHandlers(server)

	// object
	server.Echo.Use(middleware.Logger())
	server.Echo.GET("/:bucketName/:objectName", objectHanlder.GetObject)
	server.Echo.PUT("/:bucketName/:objectName", objectHanlder.PutObject)
	server.Echo.DELETE("/:bucketName/:objectName", objectHanlder.DeleteObject)

	// bucket
	server.Echo.PUT("/:bucketName", bucketHandler.CreateBucket)
}
