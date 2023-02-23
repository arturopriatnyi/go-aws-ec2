package http

//go:generate mockgen -source=handler.go -destination=mock.go -package=http

import (
	"encoding/json"
	"net/http"
	"time"

	"go-aws-ec2/pkg/counter"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type CounterManager interface {
	Add(id string) error
	Get(id string) (counter.Counter, error)
	Inc(id string) error
	Delete(id string) error
}

func NewHandler(l *zap.Logger, cm CounterManager) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), withInternalServerErrorCounter())
	r.HandleMethodNotAllowed = true

	r.NoRoute(noRoute())
	r.NoMethod(noMethod())
	r.GET("/health", health())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	counters := r.Group("/counters")
	counters.POST("", addCounter(l, cm))
	counters.GET("/:id", getCounter(l, cm))
	counters.GET("/:id/inc", incCounter(l, cm))
	counters.DELETE("/:id", deleteCounter(l, cm))

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
		start := time.Now()

		var r addCounterRequest
		if err := c.BindJSON(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := cm.Add(r.ID)

		switch err {
		case nil:
			c.AbortWithStatus(http.StatusCreated)

			defer addCounterRequestDurationHistogram.With(nil).Observe(time.Since(start).Seconds())
			defer countersNumberGauge.With(nil).Inc()
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

type getCounterResponse struct {
	ID    string `json:"id"`
	Value uint64 `json:"value"`
}

func getCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		c, err := cm.Get(id)

		switch err {
		case nil:
			ctx.AbortWithStatusJSON(http.StatusOK, getCounterResponse{
				ID:    c.ID,
				Value: c.Value,
			})
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func incCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := cm.Inc(id)

		switch err {
		case nil:
			ctx.AbortWithStatus(http.StatusOK)

			defer incCounterCounter.With(nil).Inc()
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func deleteCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := cm.Delete(id)

		switch err {
		case nil:
			ctx.AbortWithStatus(http.StatusNoContent)

			defer countersNumberGauge.With(nil).Dec()
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
