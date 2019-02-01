package restful

import (
	"advanced-database-course-project-server/config"
	"advanced-database-course-project-server/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HealthCheck(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := writer.Write([]byte("don't worry, this server is healthy!"))
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to write http response")
	}

	request.Close = true
}

func NewHttpServer() error {
	// Get apis
	router := httprouter.New()
	router.GET("/healthcheck", HealthCheck)

	return http.ListenAndServe(config.Config.RestApiConfig.Host, router)
}
