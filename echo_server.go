package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
	"os"
	"fmt"
	"net/http"
)

func skipper(c echo.Context) bool {
	return false
}

func setup_api_version_on_router(
	api_version string, server *echo.Echo,
	logger_config middleware.LoggerConfig) *echo.Echo {

	api_subrouter := server.Group(api_version)

	api_subrouter.Get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "using v1 API")
	})

	api_subrouter.Use(middleware.LoggerWithConfig(logger_config))
	api_subrouter.Use(middleware.Recover())
	return server

}


func run_simple_server(host string, port int, logger_config middleware.LoggerConfig) {
	server := setup_api_version_on_router("/v1", echo.New(), logger_config)
	server.Run(standard.New(fmt.Sprintf("%v:%v", host, port)))
}

func main() {
	CustomLoggerConfig := middleware.LoggerConfig{
		Skipper: skipper,
		Format: `${time_rfc3339}" | from: "${remote_ip}" | ` +
			`HTTP method: "${method}" | Response code: ${status} | ` +
			`Request latency: ${latency_human} | ` +
			`Request data size (in bytes): ${bytes_in} | ` +
			`Response data size (in bytes): ${bytes_out} |` + "\n",
		Output: os.Stdout,
	}

	run_simple_server("localhost", 8000, CustomLoggerConfig)
}
