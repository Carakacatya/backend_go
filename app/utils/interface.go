package utils

import "praktikum3/app/model"

type PasswordChecker interface {
    Check(hash, password string) bool
}

type TokenGenerator interface {
    Generate(user model.User) (string, error)
}

// implementasi asli
type RealPasswordChecker struct{}

func (RealPasswordChecker) Check(hash, password string) bool {
    return CheckPassword(hash, password)
}

type RealTokenGenerator struct{}

func (RealTokenGenerator) Generate(user model.User) (string, error) {
    return GenerateToken(user)
}
