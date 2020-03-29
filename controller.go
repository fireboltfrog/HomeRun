package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type controllers struct {
	Period uint `yaml:"period"`
	Margin uint `yaml:"margin"`
	Min    uint `yaml:"min"`
	Max    uint `yaml:"max"`
	Nav    uint `yaml:"nav"`
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
	back := rpio.Pin(c.Nav)
	down := rpio.Pin(c.Down)
	hold := rpio.Pin(c.Hold)
	up := rpio.Pin(c.Up)
	for _, p := range []rpio.Pin{back, down, hold, up} {
		p.Output()
		p.High()
	}
	last := time.Now()
	period := time.Second * time.Duration(c.Period)
	margin := time.Second * time.Duration(c.Margin)
	fmt.Println("wait one period")
	time.Sleep(period)
	fmt.Println("ready")
	// we want to iterate backwards, as we also navigate backwards on the remote
	// this reduces the amounts of shifts
	sort.Slice(devs, func(i, j int) bool {
		return devs[i].Device.Position > devs[j].Device.Position
	})
	var active uint = 0
	for {
		for _, d := range devs {
			if d.changed() {
				fmt.Printf("TRIGGERED: %v\n", d.Device.Position)
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
					toggle(back)
					active += c.Max - 1
					active = active % c.Max
					fmt.Printf("shifted to %v\n", active+1)
				}
				t := d.target()
				c := d.current()
				fmt.Printf("target: %v\n", t)
				switch t {
					case 0:
						toggle(down)
					case 100:
						toggle(up)
					default:
						diff := t - c
						if diff < 0 {
							toggle(down)
							diff = -diff
						} else {
							toggle(up)
						}
						duration := d.Device.Time / 100 * diff
						time.Sleep(time.Millisecond * time.Duration(duration))
						toggle(hold)

				}
				last = time.Now()
				d.update()
			}
		}
		time.Sleep(time.Second)
	}
}
