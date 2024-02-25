package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags users
// @Accept  json
// @Produce  plain
// @Param user body request.Login true "User to login"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 401 {object} response.Error "Invalid password"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var requestUser request.Login

	fmt.Println("endpoint POST /users/login called")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestUser); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload "+err.Error())
		return
	}

	if err := requestUser.Validate(); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := application.GetUserByEmail(requestUser.Email)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	if err = user.CheckPassword(requestUser.Password); err != nil {
		tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := domain.CreateToken(user)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Error while generating token")
		return
	}

	tools.FormatStrResponseBody(w, http.StatusOK, token)
}
