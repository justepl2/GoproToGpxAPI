package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags users
// @Produce  plain
// @Security BearerAuth
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint POST /users/logout called")

	accessToken := r.Context().Value("token").(string)
	if accessToken == "" {
		http.Error(w, "Access token is required", http.StatusBadRequest)
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

	// Call the GlobalSignOut function
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}
	_, err = cognito.GlobalSignOut(input)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	tools.FormatResponseBody(w, http.StatusOK, "Logout successful")
}
