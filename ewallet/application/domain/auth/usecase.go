package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"simple-go/application/entity"
	"simple-go/helper"
	"simple-go/pkg/response"
	"strings"
	"time"
)

type Repository interface {
	GetCustomerByEmail(ctx context.Context, email string) (entity.Customer, error)
	CreateCustomer(ctx context.Context, req entity.Customer) (entity.Customer, error)

	GetCustomerAuthByCustomerId(ctx context.Context, customerId int) (entity.CustomerAuth, error)
	CreateCustomerAuth(ctx context.Context, req entity.CustomerAuth) (entity.CustomerAuth, error)
	UpdateCustomerAuth(ctx context.Context, req entity.CustomerAuth, id int) (entity.CustomerAuth, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) service {
	return service{
		repository: repo,
	}
}

func (s service) RegisterCustomer(ctx context.Context, req RegisterUserRequest) response.ErrorResponse {
	emailCheck := helper.IsEmailValid(req.Email)
	if !emailCheck {
		return *response.Error("22101").WithError("email doesnt valid").WithStatusCode(http.StatusBadRequest)
	}

	newCustomerData := entity.Customer{
		Email:             req.Email,
		FullName:          req.FullName,
		BirthDate:         req.BirthDate,
		BirthPlace:        req.BirthPlace,
		IdentityPhotoLink: req.IdentityPhotoLink,
		Code:              strings.ToUpper(req.FullName[:2] + helper.SecureRandomString(5)),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	customer, err := s.repository.CreateCustomer(ctx, newCustomerData)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	newCustomerAuth := entity.CustomerAuth{
		CustomerId: customer.ID,
		Password:   string(hashedPassword),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = s.repository.CreateCustomerAuth(ctx, newCustomerAuth)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return *response.NotError()
}

func (s service) LoginCustomer(ctx context.Context, req LoginCustomerReq) (string, response.ErrorResponse) {
	emailCheck := helper.IsEmailValid(req.Email)
	if !emailCheck {
		return "", *response.Error("22101").WithError("email doesnt valid").WithStatusCode(http.StatusBadRequest)
	}

	customer, err := s.repository.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	if customer.IsEmpty {
		return "", *response.Error("22149").WithStatusCode(http.StatusBadRequest)
	}

	customerAuth, err := s.repository.GetCustomerAuthByCustomerId(ctx, customer.ID)
	if err != nil {
		return "", *response.Error("22149").WithStatusCode(http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customerAuth.Password), []byte(req.Password)); err != nil {
		return "", *response.Error("22101").WithError("Password is wrong").WithStatusCode(http.StatusBadRequest)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &entity.Claims{
		CustomerId: customer.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(helper.GetJWTKey()))
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	customerAuth.Token = tokenString
	_, err = s.repository.UpdateCustomerAuth(ctx, customerAuth, customerAuth.ID)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return tokenString, *response.NotError()
}
