package main

import (
	"net/http"

	"github.com/ktbartholomew/nsmaker/pkg/kubernetes"

	"github.com/ktbartholomew/nsmaker/pkg/types"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.Logger())
	e.Use(middleware.CSRF())

	e.Static("/", "public")
	e.POST("/create", func(c echo.Context) error {
		req := &types.CreateNamespaceRequest{}
		err := c.Bind(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &types.ErrorResponse{Message: err.Error()})
		}

		if req.Namespace == "" || req.Username == "" {
			return c.JSON(http.StatusBadRequest, &types.ErrorResponse{
				Message: "neither namespace or username can be empty",
			})
		}

		_, err = kubernetes.CreateNamespaceForUser(req.Namespace, req.Username)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &types.ErrorResponse{Message: err.Error()})
		}

		return c.NoContent(201)
	})

	e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}
