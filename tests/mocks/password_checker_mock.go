package mocks

type PasswordCheckerMock struct {
	CheckFunc func(hash, password string) bool
}

func (m PasswordCheckerMock) Check(hash, password string) bool {
	return m.CheckFunc(hash, password)
}
