package main

import (
	"fmt"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type WindowCovering struct {
	*accessory.Accessory
	WindowCovering *service.WindowCovering
}

func NewWindowCovering(info accessory.Info) *WindowCovering {
	acc := WindowCovering{}
	acc.Accessory = accessory.New(info, accessory.TypeWindowCovering)
	acc.WindowCovering = service.NewWindowCovering()
	acc.AddService(acc.WindowCovering.Service)
	return &acc
}

// Update will update the position if the current position does not match the target position
func (a *WindowCovering) Update(ch chan<- command) {
	c := a.WindowCovering.CurrentPosition.GetValue()
	t := a.WindowCovering.TargetPosition.GetValue()
	if t != c {
		fmt.Printf("update current position: %v -> %v\n", c, t)
		ch <- command{
			DeviceID: a.ID,
			Current:  c,
			Target:   t,
		}
		a.WindowCovering.CurrentPosition.SetValue(t)
	}
}
