package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct{
	Number int64 `json:"number"`
	Token string `json:"token"`
}

type LoginRequest struct{
	Number int64 `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct{
	ToAccount int `json:"toAccount"`
	Amount int `json:"amount"`
}

type CreateAccountRequest struct{
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}

type Account struct {
	ID        int `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EncryptedPassword  string `json:"-"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
	CreatedTime time.Time `json:"createdAt"`
}

func (a *Account) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword),[]byte(password)) == nil
}

func NewAccount(firstName, LastName string, password string) (*Account,error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return nil ,err
	}
	
	return &Account{
		FirstName: firstName,
		LastName: LastName,
		EncryptedPassword: string(encryptedPass),
		Number: int64(rand.Intn(999999)),
		CreatedTime: time.Now().UTC() ,
	},nil
}