package http

//go:generate mockgen -source=handler.go -destination=mock.go -package=http

import (
	"encoding/json"
	"net/http"

	"go-aws-ec2/pkg/counter"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CounterManager interface {
	Add(id string) error
}

func NewHandler(l *zap.Logger, cm CounterManager) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.HandleMethodNotAllowed = true

	r.NoRoute(noRoute())
	r.NoMethod(noMethod())
	r.GET("/health", health())

	counters := r.Group("/counters")
	counters.POST("", addCounter(l, cm))

	return r
}

func noRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func noMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	}
}

type addCounterRequest struct {
	ID string `json:"id"`
}

func addCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r addCounterRequest
		if err := c.BindJSON(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := cm.Add(r.ID)

		switch err {
		case nil:
			c.AbortWithStatus(http.StatusCreated)
		case counter.ErrExists:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			body, _ := json.Marshal(r)
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.String("body", string(body)),
				zap.Error(err),
			)

			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
