package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.GP15

	// GP14/GP15 are driven by PWM slice 7 on the RP2350/RP2040.
	pwm := machine.PWM7
	err := pwm.Configure(machine.PWMConfig{
		// RP2040/RP2350 PWM uses a 16-bit counter + max ~256x divider,
		// so the slowest achievable rate is ~7 Hz (~135 ms period). A
		// 2-second period cannot be generated directly, so we run a 200 Hz
		// carrier and toggle the duty in software to blink slowly.
		Period: uint64(5 * time.Millisecond),
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
		for i := 1; i < 255; i++ {
			pwm.Set(ch, pwm.Top()/(uint32(i)))
			time.Sleep(5 * time.Millisecond)
		}
	}
}
