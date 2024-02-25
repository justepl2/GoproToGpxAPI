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

// ForgotPassword godoc
// @Summary Send a password reset email
// @Description Send a password reset email
// @Tags users
// @Accept  json
// @Produce  plain
// @Param user body request.ForgotPassword true "User to send password reset email"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/forgot_password [post]
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var forgotPasswordRequest request.ForgotPassword

	fmt.Println("endpoint POST /users/forgot_password called")

	err := json.NewDecoder(r.Body).Decode(&forgotPasswordRequest)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = forgotPasswordRequest.Validate()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Failed to create AWS session")
		return
	}

	// Create a new Cognito client
	cognito := cognitoidentityprovider.New(sess)

	// Call the ForgotPassword function
	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
		Username: aws.String(forgotPasswordRequest.Email),
	}
	_, err = cognito.ForgotPassword(input)
	if err != nil {
		http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
		return
	}

	tools.FormatResponseBody(w, http.StatusOK, "Password reset email sent")
}
