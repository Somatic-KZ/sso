package domain

import (
	"strings"
)

const (
	TokenStatusNew              = "new"
	TokenStatusOnCheck          = "on_check"
	TokenStatusValidationFailed = "validation_failed"
	TokenStatusFinish           = "finish"
	PhoneTokenLen               = 4
)

func NewTokenizer(isTesting bool) *Tokenizer {
	return &Tokenizer{
		isTesting: isTesting,
	}
}

type Tokenizer struct {
	isTesting bool
}

func (t Tokenizer) NewTokenForPhone() (string, error) {
	if t.isTesting {
		return strings.Repeat("1", PhoneTokenLen), nil
	}

	return GenerateRandomNumbers(PhoneTokenLen)
}
