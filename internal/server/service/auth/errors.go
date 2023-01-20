package auth

import "errors"

// Ошибки возвращаемые пакетом

var TokenNotFound = errors.New("TokenNotFound") // 404
var AuthenticationError = errors.New("AuthenticationError")  // 401
var LoginAlreadyExist = errors.New("LoginAlreadyExist") // 409
var InternalServerError = errors.New("InternalServerError") // 500

