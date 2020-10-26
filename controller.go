package AGScheduler

import (
	"context"
	"time"
)

type Controller struct {
	Ctx      context.Context
	Deadline context.Context
	Cancel   context.CancelFunc
}

func NewController() *Controller {
	ctx := context.Background()
	deadline, cancel := context.WithDeadline(ctx, EmptyDateTime)
	return &Controller{
		Ctx:      ctx,
		Deadline: deadline,
		Cancel:   cancel,
	}
}

func (c *Controller) Reset(deadlineTime time.Time) {
	c.Deadline, c.Cancel = context.WithDeadline(c.Ctx, deadlineTime)
}
