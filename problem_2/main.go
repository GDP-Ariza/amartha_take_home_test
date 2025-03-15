package main

import (
	"fmt"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"
const DATE_FORMAT = "2006-01-02"

func main() {
	transactionFilePath := "transactions.csv"
	bankStatementPaths := []string{
		"bca_statements.csv",
		"mandiri_statements.csv",
	}
	startDate := "2025-01-03"
	endDate := "2025-01-03"
	err := reconciliation(transactionFilePath, bankStatementPaths, startDate, endDate)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func reconciliation(transationPath string, bankStatementPaths []string, startDate, endDate string) error {
	start, err := time.Parse(DATE_FORMAT, startDate)
	if err != nil {
		return err
	}
	end, err := time.Parse(DATE_FORMAT, endDate)
	if err != nil {
		return err
	}
	// By default, the time for start and end is set to 00:00:00, this could be problematic if startDate == endDate.
	// To mitigate this, we will add endDate by 1 day, and make it exclusive
	end = end.Add(24 * time.Hour)

	rawTrx, err := ReadTransactionFile(transationPath)
	if err != nil {
		return err
	}

	mapTrx := make(map[string]*Transaction, 0)
	for _, val := range rawTrx {
		if (start.Before(val.Time) || start.Equal(val.Time)) && val.Time.Before(end) {
			mapTrx[val.ID] = &val
		}
	}

	transactionProcessed := len(mapTrx)

	rawBankStatements := ReadAllBankStatements(bankStatementPaths)
	unmatchedBankStatements := make(map[string][]BankStatement, 0)
	deltaDiff := Currency(0)
	totalMatchedTrx := 0
	for _, val := range rawBankStatements {
		if (start.Before(val.Date) || start.Equal(val.Date)) && val.Date.Before(end) {
			if matchedTrx, ok := mapTrx[val.ReferenceID]; ok {
				// transaction matched
				totalMatchedTrx++
				diff := Currency(0)
				if matchedTrx.Type == DEBIT {
					// val.Amount is negative
					diff = matchedTrx.Amount + val.Amount
				} else if matchedTrx.Type == CREDIT {
					// val.Amount is positive
					diff = val.Amount - matchedTrx.Amount
				}

				if diff != 0 {
					// positive means we have more money than expected,
					// while negative means we have less money than expected
					deltaDiff += diff
				}

				delete(mapTrx, val.ReferenceID)
			} else {
				// no transaction found

				if _, ok := unmatchedBankStatements[val.BankCode]; !ok {
					unmatchedBankStatements[val.BankCode] = make([]BankStatement, 0)
				}
				unmatchedBankStatements[val.BankCode] = append(unmatchedBankStatements[val.BankCode], val)
			}
		}
	}

	unmatchedTransactions := make([]Transaction, 0)
	for _, val := range mapTrx {
		unmatchedTransactions = append(unmatchedTransactions, *val)
	}

	fmt.Println("=========================== Summary ===========================")
	fmt.Printf("Total transactions processed: %d\n", transactionProcessed)
	fmt.Printf("Total matched transactions: %d\n", totalMatchedTrx)
	fmt.Printf("Disrepancies: %s\n", deltaDiff.String())
	fmt.Printf("\n\n")
	fmt.Println("===============================================================")
	fmt.Printf("Total unmatched transactions: %d\n", len(mapTrx))
	fmt.Println("Details:")
	for _, val := range mapTrx {
		fmt.Println(val.String())
	}
	fmt.Printf("\n\n")
	fmt.Println("===============================================================")
	fmt.Printf("Total unmatched bank statements: %d\n", len(unmatchedBankStatements))
	fmt.Println("Details:")
	for _, arr := range unmatchedBankStatements {
		for _, val := range arr {
			fmt.Println(val.String())
		}
	}

	return nil
}
