package videos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/joho/godotenv"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Message string `json:"message"`
}

func TestCreate(t *testing.T) {
	loadEnv()
	video := &domain.Video{}
	video.FromRequest(request.CreateVideo{
		Name:     "videoTest",
		Duration: 6.66,
	})
	jsonVideo, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonVideo))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Create(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	expectedPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
	match, _ := regexp.MatchString(expectedPattern, resp.Message)
	assert.True(t, match)
}

func TestCreate_BadRequest(t *testing.T) {
	loadEnv()

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte("not a valid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_MissingFields(t *testing.T) {
	loadEnv()
	video := &domain.Video{}
	video.FromRequest(request.CreateVideo{
		// Name est manquant
		Duration: 6.66,
	})

	jsonVideo, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonVideo))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_EmptyName(t *testing.T) {
	loadEnv()
	video := &domain.Video{}
	video.FromRequest(request.CreateVideo{
		Name:     "",
		Duration: 6.66,
	})

	jsonVideo, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonVideo))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_NegativeDuration(t *testing.T) {
	loadEnv()
	video := &domain.Video{}
	video.FromRequest(request.CreateVideo{
		Name:     "Negative duration",
		Duration: -6.66,
	})

	jsonVideo, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonVideo))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func loadEnv() {
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}
}
