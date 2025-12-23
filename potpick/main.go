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
	ledPWM := machine.PWM7
	err := ledPWM.Configure(machine.PWMConfig{
		// 500hz
		Period: uint64(1*time.Second) / 500,
	})
	if err != nil {
		println("could not configure PWM:", err.Error())
		return
	}

	ch, err := ledPWM.Channel(led)
	if err != nil {
		println("could not get channel for pin:", err.Error())
		return
	}

	for {
		val := sensor.Get()
		// Set the duty cycle based on the value from the sensor (or pot)
		// The PWM is effectively an analog out that controls the LED brightness.
		duty := uint32(val) * ledPWM.Top() / 65535
		ledPWM.Set(ch, duty)

		println(val)
		time.Sleep(time.Millisecond * 100)
	}
}
