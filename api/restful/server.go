package restful

import (
	"advanced-database-course-project-server/config"
	"advanced-database-course-project-server/constants"
	"advanced-database-course-project-server/log"
	"advanced-database-course-project-server/model"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
	"net/http"
)

var ElasticsearchClient *elastic.Client

func HealthCheck(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := writer.Write([]byte("don't worry, this server is healthy!"))
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to write http response")
	}

	request.Close = true
}

func SearchRestaurant(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryVars := request.URL.Query()

	name := queryVars.Get("name")
	city := queryVars.Get("city")
	country := queryVars.Get("country")
	averageCost := queryVars.Get("average_cost")
	address := queryVars.Get("address")
	rate := queryVars.Get("rate")

	ctx := request.Context()

	search := ElasticsearchClient.Search().
		Index(constants.ElasticsearchRestaurantIndex).
		Type(constants.ElasticsearchRestaurantType)

	if name != "" && name != "NaN" {

		search = search.Query(elastic.NewMatchPhraseQuery("name", name))
	}

	if city != "" && city != "NaN" {
		search = search.Query(elastic.NewMatchPhraseQuery("city", city))
	}

	if country != "" && country != "NaN" {
		search = search.Query(elastic.NewMatchPhraseQuery("country", country))
	}

	if averageCost != "" && averageCost != "NaN" {
		search = search.Query(elastic.NewRangeQuery("average_cost").From(averageCost).To(999999999))
	}

	if address != "" && address != "NaN" {
		search = search.Query(elastic.NewMatchPhraseQuery("address", address))
	}

	if rate != "" && rate != "NaN" {
		search = search.Query(elastic.NewRangeQuery("rate").From(rate).To(5))
	}

	searchResult, err := search.From(0).Size(20).Do(ctx)

	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to search restaurants in elasticsearch")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	var restaurants []model.Restaurant

	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var r model.Restaurant
			err := json.Unmarshal(*hit.Source, &r)
			if err != nil {
				log.StdoutLogger.WithError(err).Error("failed to unmarshal restaurant search result")
			}

			r.DocumentId = hit.Id

			restaurants = append(restaurants, r)
		}

	} else {
		log.StdoutLogger.Info("search result is empty")
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(writer).Encode(restaurants)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to write response")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

}

func CreateRestaurant(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var restaurant model.Restaurant

	err := json.NewDecoder(request.Body).Decode(&restaurant)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("json decode failed")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	ctx := request.Context()

	indexResponse, err := ElasticsearchClient.Index().
		Index(constants.ElasticsearchRestaurantIndex).
		Type(constants.ElasticsearchRestaurantType).
		BodyJson(restaurant).
		Do(ctx)

	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to index the restaurant in elasticsearch")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	restaurant.DocumentId = indexResponse.Id

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(writer).Encode(restaurant)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to write response")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	request.Close = true
}

func UpdateRestaurant(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var restaurantMap map[string]interface{}

	err := json.NewDecoder(request.Body).Decode(&restaurantMap)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("json decode failed")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	ctx := request.Context()

	updateResponse, err := ElasticsearchClient.Update().
		Index(constants.ElasticsearchRestaurantIndex).
		Type(constants.ElasticsearchRestaurantType).
		Id(restaurantMap["document_id"].(string)).
		Doc(restaurantMap).
		Do(ctx)

	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to update the restaurant in elasticsearch")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	restaurantMap["document_id"] = updateResponse.Index

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(writer).Encode(restaurantMap)
	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to write response")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	request.Close = true
}

func DeleteRestaurant(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	documentId := params.ByName("document_id")

	ctx := request.Context()

	_, err := ElasticsearchClient.Delete().
		Index(constants.ElasticsearchRestaurantIndex).
		Type(constants.ElasticsearchRestaurantType).
		Id(documentId).
		Do(ctx)

	if err != nil {
		log.StdoutLogger.WithError(err).Error("failed to delete the restaurant in elasticsearch")
		writer.WriteHeader(http.StatusBadRequest)
		request.Close = true
		return
	}

	writer.WriteHeader(http.StatusOK)

	request.Close = true
}

func NewHttpServer() error {
	// Get apis
	router := httprouter.New()
	router.GET("/healthcheck", HealthCheck)
	router.GET("/restaurant/search", SearchRestaurant)
	router.POST("/restaurant", CreateRestaurant)
	router.PUT("/restaurant", UpdateRestaurant)
	router.DELETE("/restaurant/:document_id", DeleteRestaurant)

	return http.ListenAndServe(config.Config.RestApiConfig.Host, router)
}
