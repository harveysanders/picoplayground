package main

import (
	"errors"
	"log/slog"
	"machine"
	"strconv"
	"time"

	"github.com/harveysanders/picoplayground/mqttsensor/lcd"
	"github.com/harveysanders/picoplayground/mqttsensor/mqtt"
	"tinygo.org/x/drivers/hd44780i2c"
)

// mqttServerAddr is the address to the MQTT broker.
// It can be passed via linker flags.
//
// Ex: "10.0.0.9:1883"
// make flash/mqtt WIFI_SSID=spacecataz WIFI_PASS=foreigner
// tinygo build -ldflags="-X 'main.mqttServerAddr=10.0.0.9:1883'
var mqttServerAddr string

// mqttUsername is the MQTT broker username for authentication.
// Optional - if empty, anonymous connection is attempted.
// Can be passed via linker flags.
var mqttUsername string

// mqttPassword is the MQTT broker password for authentication.
// Optional - only used if mqttUsername is also set.
// Can be passed via linker flags.
var mqttPassword string

const (
	max16Bit uint16  = 65535 // Max ADC value. The Pico has an onboard 16-bit ADC.
	sysV     float32 = 3.3   // Logic level in volts. Pico runs at 3.3VDC.
)

func main() {
	start := time.Now()
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	debugLED := machine.GP21
	debugLED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.InitADC()
	sensor := machine.ADC{Pin: machine.ADC0}
	sensor.Configure(machine.ADCConfig{})

	// Setup LCD display over I2C
	err := machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.GP4,
		SCL: machine.GP5,
	})
	if err != nil {
		printErrForever(logger, "configure I2C", slog.Any("reason", err))
	}

	lcdDev, err := configureLCD(machine.I2C0)
	if err != nil {
		for {
			println(err.Error())
			time.Sleep(time.Second)
		}
	}

	lcdDev.ClearDisplay()

	// Create channel and start LCD handler goroutine
	lcdMessages := make(chan lcd.Message, 10)
	handler := lcd.NewHandler(lcdDev, lcdMessages, logger)
	go handler.Run()

	c := &mqtt.Client{
		ID:         "tinygo-mqtt",
		Logger:     logger,
		Timeout:    5 * time.Second,
		TCPBufSize: 2030, // MTU - ethhdr - iphdr - tcphdr
		Username:   mqttUsername,
		Password:   mqttPassword,
	}

	// Buffered channel of 10 readings. We may need to adjust depending
	// on the sensor read frequency and network availability
	sensorReadings := make(chan mqtt.SensorReading, 10)
	// Need a single to know when NTP sync is complete. We won't record messages until device time
	// is syncd with NTP (or fails)
	ntpDone := make(chan struct{})
	go func() {
		err := c.ConnectAndPublish(mqttServerAddr, sensorReadings, lcdMessages, ntpDone)
		if err != nil {
			// Print error in a loop in case the serial monitor is not
			// ready before the inital messages
			printErrForever(logger, "connect to MQTT broker", slog.Any("reason", err))
		}
	}()

	time.Sleep(10 * time.Second)

	// Read sensor, display readings on LCD and send off to MQTT broker
	// _________________________________________________________________

	// Buffer for LCD characters (16x2)
	// We need a preallocated buffer so the heap isn't exhausted
	// by many calls to fmt functions.
	line1 := make([]byte, 0, 20)
	line2 := make([]byte, 0, 20)
	const floatNoExp = 'f'
	for {
		// reslice the buffers to zero-length so append continues to work
		line1 = line1[:0]
		line2 = line2[:0]

		val := sensor.Get()
		percentage := (float32(val) / float32(max16Bit))
		voltage := percentage * sysV

		// Ex line1: "V: 3.2, 97.0%"
		line1 = append(line1, "V: "...)
		line1 = strconv.AppendFloat(line1, float64(voltage), floatNoExp, 1, 32)
		line1 = append(line1, ", "...)
		line1 = strconv.AppendFloat(line1, float64(percentage*100), floatNoExp, 1, 32)
		line1 = append(line1, "%"...)

		// Ex line2: "16-bit: 63452"
		line2 = append(line2, "16-bit: "...)
		line2 = strconv.AppendUint(line2, uint64(val), 10)

		// Non-blocking send to LCD
		select {
		case lcdMessages <- lcd.Message{Line1: line1, Line2: line2}:
		default:
			// Channel full - drop message to keep sensor loop responsive
		}

		// Non-blocking send to prevent main loop from blocking when channel is full
		reading := mqtt.SensorReading{
			Voltage:     voltage,
			RawUInt16:   val,
			SinceBootNS: time.Since(start),
		}
		// Only set Timestamp if NTP sync succeeded
		if !c.TimeSyncedAt.IsZero() {
			reading.Timestamp = time.Now()
		}

		select {
		case <-ntpDone:
			select {
			case sensorReadings <- reading:
			default:
				// Channel full - drop this reading to keep LCD responsive
			}
		default:
			// Still waiting on NTP. Don't send any readings to MQTT yet.
		}

		debugLED.High()
		time.Sleep(250 * time.Millisecond)
		debugLED.Low()
		time.Sleep(250 * time.Millisecond)
	}
}

// configureLCD takes a preconfigured I2C peripheral and attempts to
// initialize the HD44780 LCD display. If no LCD found on the commond I2C
// addresses (0x27, 0x3F), an error is returned.
func configureLCD(i2c *machine.I2C) (hd44780i2c.Device, error) {
	// Try common addresses (0x27 then 0x3F)
	addrs := []uint8{0x27, 0x3F}
	var lcd hd44780i2c.Device
	found := false

	for _, a := range addrs {
		println("checking I2C address...")
		dev := hd44780i2c.New(i2c, a)
		dev.Configure(hd44780i2c.Config{
			Width:  16,
			Height: 2,
		})

		lcd = dev
		found = true
		break
	}

	if !found {
		return lcd, errors.New("LCD not found on addresses: 0x27, 0x3f")
	}
	return lcd, nil
}

// printError prints a string to serial @ 1hz. It
// blocks forever.
func printErrForever(logger *slog.Logger, msg string, args ...any) {
	for {
		logger.Error(msg, args...)
		time.Sleep(time.Second * 10)
	}
}
