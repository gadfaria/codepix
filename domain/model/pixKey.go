package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKeyRepositoryInterface interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

const (
	PixKeyStatusActive   string = "active"
	PixKeyStatusInactive string = "inactive"
)

type PixKey struct {
	Base      `valid:"required"`
	Kind      string   `json:"kind" valid:"notnull"`
	Key       string   `json:"key" valid:"notnull"`
	Account   *Account `valid:"-"`
	AccountID string   `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Status    string   `json:"status" valid:"notnull"`
}

func (pixKey *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key")
	}

	if pixKey.Status != PixKeyStatusActive && pixKey.Status != PixKeyStatusInactive {
		return errors.New("invalid status")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewPixKey(key, kind string, account *Account) (*PixKey, error) {
	pixKey := PixKey{
		Kind:      kind,
		Key:       key,
		Account:   account,
		AccountID: account.ID,
		Status:    PixKeyStatusActive,
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}
