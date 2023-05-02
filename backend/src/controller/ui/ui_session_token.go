package ui

func (c *Controller) getSessionToken() (string, error) {
	sessionToken := c.sessionService.GetSessionToken()

	return sessionToken, nil
}
