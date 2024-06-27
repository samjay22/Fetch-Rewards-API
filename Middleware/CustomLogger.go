package Middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"time"
)

func Logger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return LoggerWithConfig(logger)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			logger.Info().Msg("Completed in " + stop.Sub(start).String() + "\n" + "host: " + req.Host + "\n")
			res.Flush()

			return next(c)
		}
	}
}
