package response

type Login struct {
	AccessToken  string `json:"accessToken" example:"eyJz9sdfsdf..."`
	RefreshToken string `json:"refreshToken" example:"eyJz9sdfsdf..."`
}
