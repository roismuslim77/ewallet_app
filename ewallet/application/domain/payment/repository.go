package payment

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

func (r repository) GetCustomerById(ctx context.Context, id int) (entity.Customer, error) {
	var data entity.Customer
	result := r.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) GetCustomerByCode(ctx context.Context, code string) (entity.Customer, error) {
	var data entity.Customer
	result := r.db.Where("code = ?", code).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) GetWalletTrx(ctx context.Context, customerId int) ([]entity.WalletTransaction, error) {
	var data []entity.WalletTransaction
	result := r.db.Where("customer_id = ?", customerId).Order("created_at desc").Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (r repository) GetWalletTrxByReferenceId(ctx context.Context, referenceId string) (entity.WalletTransaction, error) {
	var data entity.WalletTransaction
	result := r.db.Where("reference_id = ?", referenceId).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) CreateWalletTransaction(ctx context.Context, req entity.WalletTransaction) (entity.WalletTransaction, error) {
	if err := r.db.Debug().Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("trx is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) UpdateWalletTransaction(ctx context.Context, req entity.WalletTransaction, id int) (entity.WalletTransaction, error) {
	result := r.db.Clauses(&clause.Returning{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			result.Error = errors.New("trx is already exist")
			return req, result.Error
		}
		return req, result.Error
	}

	if result.RowsAffected < 1 {
		return req, errors.New("failed to update trx")
	}

	return req, nil
}

func (r repository) GetWalletCustomer(ctx context.Context, customerId int) (entity.Wallet, error) {
	var data entity.Wallet
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

func (r repository) CreateWalletCustomer(ctx context.Context, req entity.Wallet) (entity.Wallet, error) {
	if err := r.db.Debug().Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("wallet is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) UpdateWalletCustomer(ctx context.Context, req entity.Wallet, id int) (entity.Wallet, error) {
	result := r.db.Clauses(&clause.Returning{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			result.Error = errors.New("wallet is already exist")
			return req, result.Error
		}
		return req, result.Error
	}

	if result.RowsAffected < 1 {
		return req, errors.New("failed to update wallet")
	}

	return req, nil
}
