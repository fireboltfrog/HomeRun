package main

import "fmt"

func NewRemote(ch <-chan command) {
	for {
		msg := <-ch
		fmt.Println(msg)
	}
}
