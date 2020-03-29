package main

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type WindowCovering struct {
	*accessory.Accessory
	WindowCovering *service.WindowCovering
	Device         device
}

type device struct {
	Name       string `yaml:"name"`
	Controller uint   `yaml:"controller"`
	Position   uint   `yaml:"position"`
	Time       int   `yaml:"time"`
}

func NewWindowCovering(info accessory.Info, d device) *WindowCovering {
	acc := WindowCovering{}
	acc.Accessory = accessory.New(info, accessory.TypeWindowCovering)
	acc.WindowCovering = service.NewWindowCovering()
	acc.AddService(acc.WindowCovering.Service)
	acc.Device = d
	return &acc
}

func (a *WindowCovering) current() int {
	return a.WindowCovering.CurrentPosition.GetValue()
}

func (a *WindowCovering) target() int {
	return a.WindowCovering.TargetPosition.GetValue()
}

func (a *WindowCovering) changed() bool {
	return a.current() != a.target()
}

func (a *WindowCovering) update() {
	t := a.target()
	a.WindowCovering.CurrentPosition.SetValue(t)
}
