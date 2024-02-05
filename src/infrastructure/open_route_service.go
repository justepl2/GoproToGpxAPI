package infrastructure

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/models"
)

type OpenRouteService struct {
	BaseURL    string
	APIKey     string
	Terrain    string
	HTTPClient *http.Client
}

func NewOpenRouteService(terrain string) *OpenRouteService {
	var t string
	switch terrain {
	case string(request.TerrainRoad):
		t = "driving-car"
	case string(request.TerrainOffroad):
		t = "cycling-mountain"
	}

	return &OpenRouteService{
		BaseURL:    os.Getenv("ORS_BASE_URL"),
		APIKey:     os.Getenv("ORS_API_KEY"),
		Terrain:    t,
		HTTPClient: &http.Client{},
	}
}

func (ors *OpenRouteService) GetRoute(start, end [2]string) (models.OpenRouteServiceResponse, error) {

	bodyReq := []byte("{\"coordinates\":[[" + start[1] + "," + start[0] + "],[" + end[1] + "," + end[0] + "]],\"extra_info\":[\"roadaccessrestrictions\"]}")

	maxRetries, err := strconv.Atoi(os.Getenv("ORS_RETRY_COUNT"))
	if err != nil {
		return models.OpenRouteServiceResponse{}, fmt.Errorf("error while converting GPX_RETRY_COUNT to int: %w", err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", ors.BaseURL+"/directions/"+ors.Terrain+"/gpx", bytes.NewBuffer(bodyReq))
		if err != nil {
			return models.OpenRouteServiceResponse{}, fmt.Errorf("error while creating request for openRouteService: %w", err)
		}

		req.Header.Set("Content-type", "application/json")
		req.Header.Set("Authorization", ors.APIKey)

		resp, err := ors.HTTPClient.Do(req)
		if err != nil {
			return models.OpenRouteServiceResponse{}, fmt.Errorf("error while calling openRouteService: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 { // if OK, decode and return the response
			var openRouteServiceGpxStruct models.OpenRouteServiceResponse
			body, _ := io.ReadAll(resp.Body)
			content := xml.NewDecoder(bytes.NewBuffer(body))
			err := content.Decode(&openRouteServiceGpxStruct)
			if err != nil {
				return models.OpenRouteServiceResponse{}, fmt.Errorf("error while unmarshal openRouteService gpx: %w", err)
			}
			return openRouteServiceGpxStruct, nil
		} else if resp.StatusCode == 429 { // Too many request wait a minute
			fmt.Println("Too many request, waiting 1 minute, error code : ", resp.StatusCode)
			time.Sleep(time.Second * time.Duration(attempt*60))
		} else if attempt == maxRetries { // Max retry
			fmt.Println("Max retry, error code : ", resp.StatusCode)
			return models.OpenRouteServiceResponse{}, fmt.Errorf("error while calling openRouteService: %w", err)
		} else {
			fmt.Println("OpenRouteService return an error: ", resp.StatusCode, " retrying...")
			time.Sleep(time.Second * time.Duration(attempt*2)) // Exponential backoff
		}
	}

	return models.OpenRouteServiceResponse{}, nil
}
