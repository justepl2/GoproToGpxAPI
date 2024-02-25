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

// ResetPassword godoc
// @Summary Reset a user's password
// @Description Reset a user's password
// @Tags users
// @Accept  json
// @Produce  plain
// @Param user body request.ResetPassword true "User to reset password"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/reset_password [post]
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPasswordRequest request.ResetPassword

	fmt.Println("endpoint POST /users/reset_password called")
	// Parse the request body to get the email, code and new password
	err := json.NewDecoder(r.Body).Decode(&resetPasswordRequest)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = resetPasswordRequest.Validate()
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

	// Call the ConfirmForgotPassword function
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
		Username:         aws.String(resetPasswordRequest.Email),
		ConfirmationCode: aws.String(resetPasswordRequest.ConfirmationKey),
		Password:         aws.String(resetPasswordRequest.NewPassword),
	}
	_, err = cognito.ConfirmForgotPassword(input)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Failed to reset password"+err.Error())
		return
	}

	tools.FormatResponseBody(w, http.StatusOK, "Password reset successfully")
}
