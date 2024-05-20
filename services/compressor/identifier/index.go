package identifier

import (
	"compressor/external/redis_service"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type request struct {
	Url                  string `json:"url"`
	LocalPathUnOptimized string `json:"localPathUnOptimized"`
	LocalPathOptimized   string `json:"localPathOptimized"`
	Status               string `json:"status"`
}

func SetStatus(id string, status string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"status": status,
	}
	return updateRequest(id, fieldsToUpdate)
}

func SetLocalPathUnOptimized(id string, path string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathUnOptimized": path,
	}
	return updateRequest(id, fieldsToUpdate)
}

func SetLocalPathOptimized(id string, path string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathOptimized": path,
	}
	return updateRequest(id, fieldsToUpdate)
}

func NewRequest(url string) (string, error) {
	redisService := redis_service.GetRedisService()
	newId := uuid.New().String()
	data := &request{
		Url:    url,
		Status: "requested",
	}

	// Serialize the struct to JSON
	dataJson, _ := json.Marshal(data)
	// Save the data for the image id
	err := redisService.Set(context.Background(), newId, dataJson, 0).Err()

	// Save the mapping from url to id
	err = redisService.Set(context.Background(), url, newId, 0).Err()
	return newId, err
}

func updateRequest(id string, fields map[string]interface{}) (*request, error) {
	redisService := redis_service.GetRedisService()
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

func GetRequest(id string) (*request, error) {
	redisService := redis_service.GetRedisService()
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

func GetRequestByUrl(url string) (*request, error) {
	redisService := redis_service.GetRedisService()
	id, err := redisService.Get(context.Background(), url).Result()
	if err != nil {
		return nil, err
	}

	return GetRequest(id)
}
