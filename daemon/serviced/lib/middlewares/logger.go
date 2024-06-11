package middlewares

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/lib/context"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"gitlab.com/sibsfps/spc/spc-1/util/uuid"
)

// LoggerMiddleware provides some extra state to the logger middleware
type LoggerMiddleware struct {
	log logging.Logger
}

// MakeLogger initializes the logger middleware function
func MakeLogger(log logging.Logger) echo.MiddlewareFunc {
	logger := LoggerMiddleware{
		log: log,
	}

	return logger.handler
}

// Logger is an echo middleware to add log to the API
func (logger *LoggerMiddleware) handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		start := time.Now()

		id := uuid.NewV7()
		ctx.Set(context.RequestID, id)

		log := logger.log.With(context.RequestID, id)
		ctx.SetLogger(log.MakeEchoLogger())

		// Get a reference to the response object.
		res := ctx.Response()
		req := ctx.Request()

		// Propagate the error if the next middleware has a problem
		if err = next(ctx); err != nil {
			ctx.Error(err)
		}

		log.WithFields(
			logging.Fields{
				context.RemoteHost:       req.RemoteAddr,
				context.RemoteUser:       "-",
				context.RequestTimestamp: start,
				context.RequestMethod:    req.Method,
				context.RequestURI:       req.RequestURI,
				context.RequestProto:     req.Proto,
				context.ResponseStatus:   res.Status,
				context.ResponseSize:     strconv.FormatInt(res.Size, 10),
				context.RequestUserAgent: req.UserAgent(),
				context.Elapsed:          time.Since(start).Nanoseconds(),
			},
		).Info()

		return
	}
}
