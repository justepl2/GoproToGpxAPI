package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

// Refresh godoc
// @Summary Refresh a token
// @Description Refresh a token
// @Tags users
// @Accept  json
// @Produce  plain
// @Param user body request.RefreshToken true "Refresh token"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/refresh [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var requestRefreshToken request.RefreshToken

	fmt.Println("endpoint POST /users/refresh called")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestRefreshToken); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload "+err.Error())
		return
	}

	if err := requestRefreshToken.Validate(); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}

	// Create a new Cognito client
	cognito := cognitoidentityprovider.New(sess)

	// Call the InitiateAuth function to refresh the token
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(requestRefreshToken.RefreshToken),
		},
		ClientId: aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
	}
	output, err := cognito.InitiateAuth(input)
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	// Send the new access token in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output.AuthenticationResult.AccessToken)
}
