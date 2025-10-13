package mocks

import "errors"

type MockHasher struct {
	ShouldFail bool
}

func (m MockHasher) Hash(password string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("simulated hash error")
	}
	return password, nil
}

func (m MockHasher) Compare(hash, password string) error {
	if m.ShouldFail {
		return errors.New("simulated compare fail")
	}
	return nil
}
