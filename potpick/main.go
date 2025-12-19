package main

import (
	"machine"
	"time"
)

func main() {
	machine.InitADC()
	led := machine.GP15
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	sensor := machine.ADC{Pin: machine.ADC0}
	sensor.Configure(machine.ADCConfig{})

	// https://tinygo.org/docs/reference/microcontrollers/pico2-w/
	pwm := machine.PWM7
	err := pwm.Configure(machine.PWMConfig{
		// 500hz
		Period: uint64(1*time.Second) / 500,
	})
	if err != nil {
		println("could not configure PWM:", err.Error())
		return
	}

	ch, err := pwm.Channel(led)
	if err != nil {
		println("could not get channel for pin:", err.Error())
		return
	}

	for {
		val := sensor.Get()
		duty := uint32(val) * pwm.Top() / 65535
		pwm.Set(ch, duty)
		println(val)
		time.Sleep(time.Millisecond * 100)
	}
}
