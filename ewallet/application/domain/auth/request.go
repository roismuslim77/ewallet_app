package auth

import (
	"time"
)

type RegisterUserRequest struct {
	Email             string    `json:"email" binding:"required"`
	FullName          string    `json:"full_name" binding:"required"`
	BirthPlace        string    `json:"birth_place" binding:"required"`
	BirthDate         time.Time `json:"birth_date" binding:"required"`
	IdentityPhotoLink string    `json:"identity_photo_link" binding:"required"`
	Password          string    `json:"password" binding:"required"`
}

type LoginCustomerReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CustomerAuthHeaderReq struct {
	Authorization string `header:"Authorization" binding:"required"`
}
