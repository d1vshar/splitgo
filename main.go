package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/d1vshar/splitgo/api"
	"github.com/d1vshar/splitgo/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var baseG = e.Group("/api/v1")

	App := core.NewApp()
	Api := &api.Api{
		App: App,
	}

	e.Use(middleware.RateLimiterWithConfig(RateLimiterConfig(Api)))

	Api.RegisterCategoryRoutes(baseG)
	Api.RegisterTransactionRoutes(baseG)

	s := http.Server{
		Addr:    ":1323",
		Handler: e,
	}

	App.Logger().Info("starting server on port 1323")

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		App.Logger().Fatal(err)
	}
}

func RateLimiterConfig(Api *api.Api) middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 10, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return api.ErrorResponse(context, http.StatusForbidden, errors.New("user not found"))
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return api.ErrorResponse(context, http.StatusTooManyRequests, errors.New("too many requests"))
		},
	}

}
