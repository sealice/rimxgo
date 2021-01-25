package controllers

import (
	"fmt"
	"time"
)

type DefaultController struct {
	baseController
}

func (c *DefaultController) Get() string {
	visits := 1
	since := time.Since(c.StartTime).Seconds()
	logger.Debugf("login user: %#v", c.User)

	// write the current, updated visits.
	if c.Session != nil {
		logger.Infof("sessionid: %s", c.Session.ID())
		visits = c.Session.Increment("visits", 1)
	}

	return fmt.Sprintf("%d visit(s) from my current session in %0.1f seconds of server's up-time",
		visits, since)
}
