package sourcerule

import "errors"

var ErrSourceRuleNotFound = errors.New("source rule not found")
var ErrValidation = errors.New("source rule validation failed")
