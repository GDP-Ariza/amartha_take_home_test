package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Currency int64

func (c *Currency) String() string {
	return fmt.Sprintf("%.2f", float64(*c)/100)
}

func NewCurrency(str string) (Currency, error) {
	arr := strings.Split(str, ".")
	res, err := strconv.ParseInt(arr[0], 10, 64)
	res *= 100
	if len(arr) == 2 {
		cents, err := strconv.Atoi(arr[1])
		if err == nil {
			if res >= 0 {
				res += int64(cents)
			} else {
				res -= int64(cents)
			}
		}
	}

	return Currency(res), err
}

type TransactionType string

const (
	DEBIT  TransactionType = "debit"
	CREDIT TransactionType = "credit"
)

type Transaction struct {
	ID     string
	Amount Currency
	Type   TransactionType
	Time   time.Time
}

func (t *Transaction) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", t.ID, t.Type, t.Time.Format(TIME_FORMAT), t.Amount.String())
}

type BankStatement struct {
	BankCode    string
	ReferenceID string
	Amount      Currency
	Date        time.Time
}

func (b *BankStatement) String() string {
	return fmt.Sprintf("%s %s, %s, %s", b.BankCode, b.ReferenceID, b.Date.Format(TIME_FORMAT), b.Amount.String())
}
