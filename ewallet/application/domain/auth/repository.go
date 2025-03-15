package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple-go/application/entity"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) GetCustomerByEmail(ctx context.Context, email string) (entity.Customer, error) {
	var data entity.Customer
	result := r.db.Where("email = ?", email).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) CreateCustomer(ctx context.Context, req entity.Customer) (entity.Customer, error) {
	if err := r.db.Clauses(clause.Returning{}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("user is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) GetCustomerAuthByCustomerId(ctx context.Context, customerId int) (entity.CustomerAuth, error) {
	var data entity.CustomerAuth
	result := r.db.Where("customer_id = ?", customerId).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) CreateCustomerAuth(ctx context.Context, req entity.CustomerAuth) (entity.CustomerAuth, error) {
	if err := r.db.Clauses(clause.Returning{}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("user is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) UpdateCustomerAuth(ctx context.Context, req entity.CustomerAuth, id int) (entity.CustomerAuth, error) {
	result := r.db.Clauses(&clause.Returning{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			result.Error = errors.New("auth is already exist")
			return req, result.Error
		}
		return req, result.Error
	}

	if result.RowsAffected < 1 {
		return req, errors.New("failed to update auth")
	}

	return req, nil
}
