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

// ValidateMail godoc
// @Summary Validate mail
// @Description Validate mail
// @Tags users
// @Accept  json
// @Produce  plain
// @Param mailValidation body request.MailValidation true "Mail validation"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response.Error "Invalid request"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /users/validateMail [post]
func ValidateMail(w http.ResponseWriter, r *http.Request) {
	var requestMailValidation request.MailValidation

	fmt.Println("endpoint POST /users/validateMail called")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestMailValidation); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := requestMailValidation.Validate(); err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
		Username:         aws.String(requestMailValidation.Username),
		ConfirmationCode: aws.String(requestMailValidation.ConfirmationCode),
	}

	confirmSignupOutput, err := svc.ConfirmSignUp(input)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = confirmSignupOutput
	tools.FormatEmptyResponseBody(w, http.StatusOK)
}
