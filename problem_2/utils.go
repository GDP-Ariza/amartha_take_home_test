package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func readCSVFile(path string) (lines [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read line by line
	for {
		row, err := reader.Read()
		if err != nil {
			break // Reaches EOF
		}

		lines = append(lines, row)
	}
	return lines, nil
}

func ReadTransactionFile(path string) (data []Transaction, err error) {
	lines, err := readCSVFile(path)
	if err != nil {
		return data, err
	}

	for i, line := range lines {
		// 1st line is header, skip
		if i == 0 {
			continue
		}

		if len(line) != 4 {
			return data, fmt.Errorf("invalid data on line %d", i)
		}

		amount, err := NewCurrency(line[1])
		if err != nil {
			return data, fmt.Errorf("invalid amount on line %d", i)
		}

		dateTime, err := time.Parse(TIME_FORMAT, line[3])
		if err != nil {
			return data, fmt.Errorf("time not in valid format: %s", line[3])
		}

		data = append(data, Transaction{
			ID:     line[0],
			Amount: amount,
			Type:   TransactionType(strings.ToLower(line[2])),
			Time:   dateTime,
		})
	}

	return data, nil
}

func readBankStatementFile(path string, bankCode string) (data []BankStatement, err error) {
	lines, err := readCSVFile(path)
	if err != nil {
		return data, err
	}

	for i, line := range lines {
		// 1st line is header, skip
		if i == 0 {
			continue
		}

		if len(line) != 3 {
			return data, fmt.Errorf("invalid data on line %d", i)
		}

		amount, err := NewCurrency(line[1])
		if err != nil {
			return data, fmt.Errorf("invalid amount on line %d", i)
		}

		dateTime, err := time.Parse(TIME_FORMAT, line[2])
		if err != nil {
			return data, fmt.Errorf("time not in valid format: %s", line[3])
		}

		data = append(data, BankStatement{
			BankCode:    bankCode,
			ReferenceID: line[0],
			Amount:      amount,
			Date:        dateTime,
		})

	}

	return data, nil
}

func ReadAllBankStatements(paths []string) (data []BankStatement) {
	for _, val := range paths {
		arr := strings.Split(val, "/")
		bankCode := strings.Split(arr[len(arr)-1], "_")[0]
		temp, err := readBankStatementFile(val, bankCode)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		data = append(data, temp...)
	}
	return data
}
