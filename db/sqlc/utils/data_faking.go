package utils

import (
	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc"
	"github.com/bxcodec/faker/v3"
)

// im going to fake data for Owner, Currency and Balance

type RandomData struct {
	Owner    string  `faker:"name"`
	Balance  float64 `faker:"amount"`
	Currency string  `faker:"currency"`
}

func To_account(randomData *RandomData) sqlc.CreateAccountParams {
	return sqlc.CreateAccountParams{
		Owner:    randomData.Owner,
		Balance:  int64(randomData.Balance),
		Currency: randomData.Currency,
	}
}

func GenerateRandomOwner() *RandomData {
	var randomData RandomData
	err := faker.FakeData(&randomData)
	if err != nil {
		panic(err)
	}
	return &randomData
}
