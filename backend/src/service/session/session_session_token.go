package session

import "fmt"

func (s *Service) GetSessionToken() string {
	return s.sessionToken
}

func (s *Service) VerifySessionToken(sessionToken string) error {
	if sessionToken == "" {
		return fmt.Errorf("empty session token")
	}

	if sessionToken != s.sessionToken {
		return fmt.Errorf("wrong session token")
	}

	return nil
}
