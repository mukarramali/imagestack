package identifier

import (
	"compressor/external/redis_service"
	"context"
	"encoding/json"
	"errors"
	"imagestack/lib"

	"github.com/google/uuid"
)

type request struct {
	Url                  string `json:"url"`
	Quality              int    `json:"quality"`
	Width                int    `json:"width"`
	LocalPathUnOptimized string `json:"localPathUnOptimized"`
	LocalPathOptimized   string `json:"localPathOptimized"`
	Status               string `json:"status"`
}

func SetStatus(id string, status string, redisUrl string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"status": status,
	}
	return updateRequest(id, fieldsToUpdate, redisUrl)
}

func SetLocalPathUnOptimized(id string, path string, redisUrl string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathUnOptimized": path,
	}
	return updateRequest(id, fieldsToUpdate, redisUrl)
}

func SetLocalPathOptimized(id string, path string, redisUrl string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathOptimized": path,
	}
	return updateRequest(id, fieldsToUpdate, redisUrl)
}

func NewRequest(url string, quality int, width int, redisUrl string) (string, error) {
	redisService := redis_service.GetRedisService(redisUrl)
	newId := uuid.New().String()
	data := &request{
		Url:     url,
		Status:  "requested",
		Quality: quality,
		Width:   width,
	}

	// Serialize the struct to JSON
	dataJson, _ := json.Marshal(data)
	// Save the data for the image id
	err := redisService.Set(context.Background(), newId, dataJson, 0).Err()
	lib.FailOnError(err, "Couldn't created request in redis")

	// Save the mapping from url to id
	err = redisService.Set(context.Background(), url, newId, 0).Err()
	return newId, err
}

func updateRequest(id string, fields map[string]interface{}, redisUrl string) (*request, error) {
	redisService := redis_service.GetRedisService(redisUrl)
	val, err := redisService.Get(context.Background(), id).Result()
	if err != nil {
		return nil, err
	}

	var data request
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	// Dynamically update fields
	if status, ok := fields["status"]; ok {
		data.Status = status.(string)
	}
	if localPathUnOptimized, ok := fields["localPathUnOptimized"]; ok {
		data.LocalPathUnOptimized = localPathUnOptimized.(string)
	}
	if localPathOptimized, ok := fields["localPathOptimized"]; ok {
		data.LocalPathOptimized = localPathOptimized.(string)
	}

	// Serialize the updated struct to JSON
	updatedDataJson, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	err = redisService.Set(context.Background(), id, updatedDataJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GetRequest(id string, redisUrl string) (*request, error) {
	redisService := redis_service.GetRedisService(redisUrl)
	val, err := redisService.Get(context.Background(), id).Result()
	if err != nil {
		return nil, err
	}

	var data request
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetRequestByUrl(url string, redisUrl string) (*request, error) {
	redisService := redis_service.GetRedisService(redisUrl)
	if redisService == nil {
		return nil, errors.New("redis service not available")
	}
	id, err := redisService.Get(context.Background(), url).Result()
	if err != nil {
		return nil, err
	}

	return GetRequest(id)
}
