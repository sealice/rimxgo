package controllers

import (
	"fmt"
	"time"
)

type DefaultController struct {
	baseController
}

func (c *DefaultController) Get() string {
	logger.Infof("sessionid: %s", c.Session.ID())
	visits := c.Session.Increment("visits", 1)
	// write the current, updated visits.
	since := time.Now().Sub(c.StartTime).Seconds()
	return fmt.Sprintf("%d visit(s) from my current session in %0.1f seconds of server's up-time",
		visits, since)
}
