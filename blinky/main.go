package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.GP15
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.High() // LED on
		println("LED on")
		time.Sleep(time.Millisecond * 1000)

		led.Low() // LED off
		println("LED off")
		time.Sleep(time.Millisecond * 1000)
	}
}
