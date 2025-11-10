package dto

import "fmt"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (l *LoginRequest) Validate() error {
	if l.Password == "" {
		return fmt.Errorf("password must not be empty")
	}

	if len(l.Password) < 6 {
		return fmt.Errorf("password must be longer than 6 symb")
	}

	if len(l.Username) < 3 {
		return fmt.Errorf("username must be longer than 3 symb")
	}

	if l.Username == "" {
		return fmt.Errorf("username must be not empty")
	}

	return nil
}
