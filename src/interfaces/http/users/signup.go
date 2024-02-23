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

// Signup godoc
// @Summary Signup a new user
// @Description Signup a new user
// @Tags users
// @Accept  json
// @Produce  plain
// @Param user body request.Signup true "User to signup"
// @Success 201 {object} response.UUIDResponse "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/signup [post]
func Signup(w http.ResponseWriter, r *http.Request) {
	var requestUser request.Signup
	var user domain.User

	fmt.Println("endpoint POST /users/signup called")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestUser); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := requestUser.Validate(); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.FromRequest(requestUser); err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot create user, err : "+err.Error())
		return
	}

	if err := application.AddUser(&user); err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot create user, err : "+err.Error())
		return
	}

	tools.FormatUuidResponseBody(w, http.StatusCreated, user.ID.String())
}
