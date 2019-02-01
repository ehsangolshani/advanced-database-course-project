package main

import (
	"advanced-database-course-project-server/api/restful"
	"advanced-database-course-project-server/config"
	"advanced-database-course-project-server/constants"
	"advanced-database-course-project-server/log"
	"context"
	"github.com/olivere/elastic"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {

	config.LoadConfiguration()
	log.InitLogrusLogger(config.Config.LogLevel, config.Config.PrettyLog)

	// maximum cpu cores value that go application uses should be more than 0 and not more than max available logical cores of underlying infrastructure
	if config.Config.GoMaxProcs > 0 && config.Config.GoMaxProcs <= runtime.NumCPU() {
		runtime.GOMAXPROCS(config.Config.GoMaxProcs)
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	elasticsearchClient, err := elastic.NewClient(elastic.SetHttpClient(&httpClient))
	if err != nil {
		log.StdoutLogger.WithError(err).Fatal("failed to initialize elasticsearch client")
	}

	restful.ElasticsearchClient = elasticsearchClient

	ctx := context.Background()

	exists, err := elasticsearchClient.IndexExists(constants.ElasticsearchRestaurantIndex).Do(ctx)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to check elasticsearch index existence")
	}

	if !exists {
		// Create a new index.
		createIndex, err := elasticsearchClient.CreateIndex(constants.ElasticsearchRestaurantIndex).Do(ctx)
		if err != nil {
			log.StdoutLogger.WithError(err).Fatal("failed to create elasticsearch index")
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// profiler endpoint
	if config.Config.ProfilerConfig.Enabled {
		go func() {
			log.StdoutLogger.Info(http.ListenAndServe(config.Config.ProfilerConfig.Host, nil))
			log.StdoutLogger.Info("profiler stopped operation!")
		}()
	}

	// rest api
	err = restful.NewHttpServer()

	if err != nil {
		log.StdoutLogger.WithError(err).Error("http server stopped!")
	}

	log.StdoutLogger.Info("application is stopping, end of story!")
}
