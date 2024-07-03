package model

import "time"

type User struct {
	UUID     string `gorm:"primaryKey;uniqueIndex"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	RegisteredAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;not null"`
	Status       bool      `gorm:"default:true;not null"`

	Accounts []Account `gorm:"foreignKey:UserUUID"`
}

type Account struct {
	UUID          string `gorm:"primaryKey;uniqueIndex"`
	UserUUID      string `gorm:"index;not null"`
	AccountNumber string `gorm:"unique;not null"`

	User    User    `gorm:"foreignKey:UserUUID"`
	Balance Balance `gorm:"foreignKey:AccountUUID"`
	//@note: account type can be added so that the currency can be known
}

type Balance struct {
	UUID        string  `gorm:"primaryKey;uniqueIndex"`
	AccountUUID string  `gorm:"index;not null"`
	Balance     float64 `gorm:"not null"`
}

type Transaction struct {
	UUID                string  `gorm:"primaryKey;uniqueIndex"`
	AccountUUID         string  `gorm:"index;not null"`
	ReceiverAccountUUID string  `gorm:"index;not null"`
	Amount              float64 `gorm:"not null"`

	Account         Account `gorm:"foreignKey:AccountUUID;references:UUID"`
	ReceiverAccount Account `gorm:"foreignKey:ReceiverAccountUUID;references:UUID"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
