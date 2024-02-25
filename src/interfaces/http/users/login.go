package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(requestUser.Email),
			"PASSWORD": aws.String(requestUser.Password),
		},
		ClientId: aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
	}

	result, err := svc.InitiateAuth(input)

	if err != nil {
		log.Println("Failed to authenticate user", err)
		return
	}

	tools.FormatStrResponseBody(w, http.StatusOK, *result.AuthenticationResult.AccessToken)
}
