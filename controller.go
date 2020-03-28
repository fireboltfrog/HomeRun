package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type controllers struct {
	Period uint `yaml:"period"`
	Margin uint `yaml:"margin"`
	Min    uint `yaml:"min"`
	Max    uint `yaml:"max"`
	Next   uint `yaml:"next"`
	Down   uint `yaml:"down"`
	Hold   uint `yaml:"hold"`
	Up     uint `yaml:"up"`
}

func toggle(pin rpio.Pin) {
	pin.Low()
	time.Sleep(200 * time.Millisecond)
	pin.High()
	time.Sleep(300 * time.Millisecond)
}

func (c *controllers) monitor(devs []*WindowCovering) {
	if err := rpio.Open(); err != nil {
		log.Fatal("coult not find GPIO pins: ", err)
	}
	next := rpio.Pin(c.Next)
	down := rpio.Pin(c.Down)
	hold := rpio.Pin(c.Hold)
	up := rpio.Pin(c.Up)
	for _, p := range []rpio.Pin{next, down, hold, up} {
		p.Output()
		p.High()
	}
	last := time.Now()
	period := time.Second * time.Duration(c.Period)
	margin := time.Second * time.Duration(c.Margin)
	fmt.Println("wait one period")
	time.Sleep(period)
	fmt.Println("ready")
	var active uint = 0
	for {
		for _, d := range devs {
			if d.changed() {
				// wake up controller if necessary
				if dur := time.Now().Sub(last); dur > period {
					if wait := period + margin - dur; 0 < wait {
						fmt.Println("wait for controller")
						time.Sleep(wait)
					}
					fmt.Println("wake up controller")
					toggle(hold)
					time.Sleep(time.Second)
				}
				for active+1 != d.Device.Position {
					toggle(next)
					active++
					active = active % c.Max
					fmt.Printf("shifted to %v\n", active+1)
				}
				c := d.current()
				t := d.target()
				fmt.Printf("target: %v\n", t)
				if t < c {
					toggle(down)
					fmt.Println("down")
				} else {
					toggle(up)
					fmt.Println("up")
				}
				last = time.Now()
				d.update()
			}
		}
		time.Sleep(time.Second)
	}
}
