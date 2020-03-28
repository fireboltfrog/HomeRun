package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"gopkg.in/yaml.v2"
)

type homekit struct {
	Pin string `yaml:"pin"`
}

type config struct {
	Devices     []device      `yaml:"devices"`
	Controllers []controllers `yaml:"controllers"`
	HomeKit     homekit       `yaml:"homekit"`
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
	rdvs := make([][]*WindowCovering, len(c.Controllers))
	for _, d := range c.Devices {
		fmt.Printf("adding %s\n", d.Name)
		info := accessory.Info{
			Name: d.Name,
		}
		ac := NewWindowCovering(info, d)
		accs = append(accs, ac.Accessory)
		rdvs[d.Controller] = append(rdvs[d.Controller], ac)
	}
	conf := hc.Config{Pin: c.HomeKit.Pin}
	bridge := accessory.NewBridge(accessory.Info{Name: "HomeRun"})
	t, err := hc.NewIPTransport(conf, bridge.Accessory, accs...)
	if err != nil {
		log.Fatal(err)
	}
	hc.OnTermination(func() {
		<-t.Stop()
	})

	for _, ds := range rdvs {
		go c.Controllers[0].monitor(ds)
	}
	fmt.Printf("HomeKit PIN: %v\n",c.HomeKit.Pin)
	t.Start()
}
