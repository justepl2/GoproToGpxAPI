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

	// if err := application.AddUser(&user); err != nil {
	// 	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot create user, err : "+err.Error())
	// 	return
	// }

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
		Password: aws.String(requestUser.Password),
		Username: aws.String(requestUser.Username),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(requestUser.Email),
			},
			{
				Name:  aws.String("given_name"),
				Value: aws.String(requestUser.FirstName),
			},
			{
				Name:  aws.String("family_name"),
				Value: aws.String(requestUser.LastName),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(requestUser.PhoneNumber),
			},
		},
	}

	_, err := svc.SignUp(input)

	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot create user, err : "+err.Error())
		return
	}

	tools.FormatEmptyResponseBody(w, http.StatusCreated)
}
