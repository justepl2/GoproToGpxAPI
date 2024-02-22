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
