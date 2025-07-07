package model

type RequestCreateUser struct {
	Username string
	Email    string
	Password string
}

type RequestLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLoginUser struct {
	AccessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}
