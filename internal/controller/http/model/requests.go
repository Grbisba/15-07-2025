package model

type (
	UploadFileRequest struct {
		URLs []string `json:"urls"`
	}
	CreateUserRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)
