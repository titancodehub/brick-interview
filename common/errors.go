package common

import "errors"

var ErrorRecordNotExist = errors.New("data not exist")
var ErrorMerchantNotExist = errors.New("merchant not exist")
var ErrorInsufficientBalance = errors.New("insufficient balance")
var ErrorDuplicateReference = errors.New("duplicate reference")

var ErrorBankAccountNotExist = errors.New("bank account not exist")
