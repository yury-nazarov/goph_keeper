package secret

import "errors"

var AuthenticationError = errors.New("AuthenticationError")
var ItemNotFound = errors.New("SecretNotFound")
var InternalServerError = errors.New("InternalServerError")
var err error
