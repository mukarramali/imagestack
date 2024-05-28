package identifier

import (
	"compressor/external/redis_service"
	"context"
	"encoding/json"
	"imagestack/lib"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	Url                  string `json:"url"`
	Quality              int    `json:"quality"`
	Width                int    `json:"width"`
	LocalPathUnOptimized string `json:"localPathUnOptimized"`
	LocalPathOptimized   string `json:"localPathOptimized"`
	Status               string `json:"status"`
}

type RequestService struct {
	redisUrl     string
	redisService *redis.Client
}

func NewRequestService(redisUrl string) *RequestService {
	return &RequestService{
		redisUrl:     redisUrl,
		redisService: redis_service.GetRedisService(redisUrl),
	}
}

func (service *RequestService) SetStatus(id string, status string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"status": status,
	}
	return service.updateRequest(id, fieldsToUpdate)
}

func (service *RequestService) SetLocalPathUnOptimized(id string, path string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathUnOptimized": path,
	}
	return service.updateRequest(id, fieldsToUpdate)
}

func (service *RequestService) SetLocalPathOptimized(id string, path string) (*request, error) {
	fieldsToUpdate := map[string]interface{}{
		"localPathOptimized": path,
	}
	return service.updateRequest(id, fieldsToUpdate)
}

func (service *RequestService) NewRequest(url string, quality int, width int) (string, error) {
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
	err := service.redisService.Set(context.Background(), newId, dataJson, 0).Err()
	lib.FailOnError(err, "Couldn't created request in redis")

	// Save the mapping from url to id
	err = service.redisService.Set(context.Background(), url, newId, 0).Err()
	return newId, err
}

func (service *RequestService) updateRequest(id string, fields map[string]interface{}) (*request, error) {
	val, err := service.redisService.Get(context.Background(), id).Result()
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

	err = service.redisService.Set(context.Background(), id, updatedDataJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (service *RequestService) GetRequest(id string) (*request, error) {
	val, err := service.redisService.Get(context.Background(), id).Result()
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

func (service *RequestService) GetRequestByUrl(url string) (*request, error) {
	id, err := service.redisService.Get(context.Background(), url).Result()
	if err != nil {
		return nil, err
	}

	return service.GetRequest(id)
}
