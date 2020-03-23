package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"gopkg.in/yaml.v2"
)

type device struct {
	Name   string `yaml:"name"`
	Remote int    `yaml:"remote"`
}

type homekit struct {
	Pin string `yaml:"pin"`
}

type remote struct {
	Sleep  uint
	Margin uint
	Count  uint
}

type config struct {
	Devices []device `yaml:"devices"`
	HomeKit homekit  `yaml:"homekit"`
	Remote  remote   `yaml:"remote"`
}

type command struct {
	DeviceID uint64
	Current  int
	Target   int
}

func main() {

	filename := flag.String("filename", "config.yml", "filename of the config file")
	flag.Parse()

	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Errorf("could not read config file: %v", err)
	}

	c := config{}
	yaml.Unmarshal(data, &c)

	accs := []*accessory.Accessory{}
	// each remote has a go routine that handles the commands
	// which are passed trough the remotes channel
	for i := 0; i < c.Remote.Count; i++ {
		rchan := make(chan command)
		// add devices assinged to this remote
		for _, d := range c.Devices {
			if d.Remote == i {
				fmt.Printf("adding %s\n", d.Name)
				info := accessory.Info{
					Name: d.Name,
				}
				ac := NewWindowCovering(info)
				accs = append(accs, ac.Accessory)
				// periodically check for updates
				go func() {
					for {
						time.Sleep(time.Second)
						ac.Update(rchan)
					}
				}()
			}
		}
		go func() {
			NewRemote(rchan)
		}()
	}

	conf := hc.Config{Pin: c.HomeKit.Pin}
	t, err := hc.NewIPTransport(conf, accs[0], accs[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
