package domain

import (
	"fmt"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
)

var TypeToPrefix = map[AccountType]string{
	AccountTypeAsset:     "AST-",
	AccountTypeLiability: "LIA-",
	AccountTypeEquity:    "EQU-",
	AccountTypeRevenue:   "REV-",
	AccountTypeExpense:   "EXP-",
}

func PrefixForType(t AccountType) (string, error) {
	prefix, ok := TypeToPrefix[t]
	if !ok {
		return "", fmt.Errorf("unknown account type: %s", t)
	}
	return prefix, nil
}
