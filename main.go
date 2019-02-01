package main

import (
	"advanced-database-course-project-server/api/restful"
	"advanced-database-course-project-server/config"
	"advanced-database-course-project-server/log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func main() {

	config.LoadConfiguration()
	log.InitLogrusLogger(config.Config.LogLevel, config.Config.PrettyLog)

	// maximum cpu cores value that go application uses should be more than 0 and not more than max available logical cores of underlying infrastructure
	if config.Config.GoMaxProcs > 0 && config.Config.GoMaxProcs <= runtime.NumCPU() {
		runtime.GOMAXPROCS(config.Config.GoMaxProcs)
	}

	// profiler endpoint
	if config.Config.ProfilerConfig.Enabled {
		go func() {
			log.StdoutLogger.Info(http.ListenAndServe(config.Config.ProfilerConfig.Host, nil))
			log.StdoutLogger.Info("profiler stopped operation!")
		}()
	}

	// rest api
	err := restful.NewHttpServer()

	if err != nil {
		log.StdoutLogger.WithError(err).Error("http server stopped!")
	}

	log.StdoutLogger.Info("application is stopping, end of story!")
}
