package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/netip"
	"runtime"
	"time"

	"github.com/harveysanders/picoplayground/mqttsensor/cyw43439"
	"github.com/harveysanders/picoplayground/mqttsensor/lcd"
	"github.com/harveysanders/picoplayground/mqttsensor/ntp"
	mqtt "github.com/soypat/natiu-mqtt"
	"github.com/soypat/seqs"
	"github.com/soypat/seqs/stacks"
)

var (
	pubFlags, _ = mqtt.NewPublishFlags(mqtt.QoS0, false, false)
	pubVar      = mqtt.VariablesPublish{
		TopicName:        []byte("tinygo-pico-test"),
		PacketIdentifier: 0xc0fe,
	}
)

type SensorReading struct {
	Voltage     float32
	RawUInt16   uint16
	SinceBootNS time.Duration
	Timestamp   time.Time // Wall-clock time. Zero if NTP sync failed.
}

type Client struct {
	ID                string
	Timeout           time.Duration
	TCPBufSize        int
	Logger            *slog.Logger
	HeartbeatInterval time.Duration
	TimeSyncedAt      time.Time // When NTP sync occurred. Zero if never synced.
	Username          string    // MQTT broker username (optional)
	Password          string    // MQTT broker password (optional, requires Username)
}

func (c *Client) ConnectAndPublish(addr string, readings <-chan SensorReading, lcdMessages chan<- lcd.Message, ntpDone chan<- struct{}) error {
	// Close the ntpDone channel if we return early on an error.
	defer func() {
		close(ntpDone)
	}()

	lcd.Send(lcdMessages, "Connecting to", "WiFi...")
	dchpClient, stack, _, err := cyw43439.SetupWithDHCP(cyw43439.SetupConfig{
		Hostname: c.ID,
		Logger:   c.Logger,
		TCPPorts: 1, // For MQTT over TCP.
		UDPPorts: 2, // For DNS (MQTT + NTP) and NTP client.
	})
	if err != nil {
		return errors.New("setup DHCP:" + err.Error())
	}

	dnsResolver, err := cyw43439.NewResolver(stack, dchpClient)
	if err != nil {
		return errors.New("dns resolver:" + err.Error())
	}

	// Get router's hardware address from resolver (it's already been resolved during DNS lookup).
	// We need the router's MAC address to send packets through the gateway to the internet.
	serverHWAddr, err := dnsResolver.RouterHWAddr()
	if err != nil {
		return errors.New("router hwaddr:" + err.Error())
	}
	c.Logger.Info(fmt.Sprintf("router hwaddr: %x", (serverHWAddr[0:6])))

	// Sync system time via NTP
	lcd.Send(lcdMessages, "Syncing time", "via NTP...")
	err = ntp.SyncTime(stack, dnsResolver, serverHWAddr, c.Logger)
	if err != nil {
		c.Logger.Error("ntp sync failed", slog.String("reason", err.Error()))
		lcd.Send(lcdMessages, "NTP sync failed", "Continuing...")
		time.Sleep(2 * time.Second)
	} else {
		c.TimeSyncedAt = time.Now()
		lcd.Send(lcdMessages, "Time synced", c.TimeSyncedAt.Format("15:04:05"))
		c.Logger.Info("ntp:success", slog.Time("time", c.TimeSyncedAt))
		time.Sleep(2 * time.Second)
	}
	// Signal we're done with NTP, even if it fails
	close(ntpDone)

	c.Logger.Info("MQTT address: " + addr)

	// Parse hostname and port from addr (e.g., "hostname:8883")
	mqttHost, portStr, err := splitHostPort(addr)
	if err != nil {
		return errors.New("parsing host:port from " + addr + ": " + err.Error())
	}

	addrs, err := dnsResolver.LookupNetIP(mqttHost)
	if err != nil {
		return errors.New("dns lookup for " + mqttHost + ": " + err.Error())
	}

	// Set up MQTT client and broker connection
	mqttAddr := addrs[0]
	c.Logger.Info("resolved IP: " + mqttAddr.String())
	port := parsePort(portStr)

	start := time.Now()
	rng := rand.New(rand.NewSource(int64(time.Now().Sub(start))))
	// Start TCP server.
	clientAddr := netip.AddrPortFrom(stack.Addr(), uint16(rng.Intn(65535-1024)+1024))
	conn, err := stacks.NewTCPConn(stack, stacks.TCPConnConfig{
		TxBufSize: uint16(c.TCPBufSize),
		RxBufSize: uint16(c.TCPBufSize),
	})
	if err != nil {
		panic("conn create:" + err.Error())
	}

	closeConn := func(err string) {
		slog.Error("tcpconn:closing", slog.String("err", err))
		conn.Close()
		for !conn.State().IsClosed() {
			slog.Info("tcpconn:waiting", slog.String("state", conn.State().String()))
			time.Sleep(1000 * time.Millisecond)
		}
	}

	cfg := mqtt.ClientConfig{
		Decoder: mqtt.DecoderNoAlloc{UserBuffer: make([]byte, 4096)},
		OnPub: func(pubHead mqtt.Header, varPub mqtt.VariablesPublish, r io.Reader) error {
			c.Logger.Info("received message", slog.String("topic", string(varPub.TopicName)))
			return nil
		},
	}
	var varconn mqtt.VariablesConnect
	varconn.SetDefaultMQTT([]byte(c.ID))

	// Set authentication credentials if provided
	if c.Username != "" {
		varconn.Username = []byte(c.Username)
		if c.Password != "" {
			varconn.Password = []byte(c.Password)
		}
	}

	mqttClient := mqtt.NewClient(cfg)

	// Connection loop for TCP+MQTT.
	for {
		random := rng.Uint32()
		c.Logger.Info("socket:listen")
		lcd.Send(lcdMessages, "addr", addr)
		err = conn.OpenDialTCP(clientAddr.Port(), serverHWAddr, netip.AddrPortFrom(mqttAddr, port), seqs.Value(random))
		if err != nil {
			panic("socket dial:" + err.Error())
		}
		lcd.Send(lcdMessages, "Connecting...", "TCP handshake")
		retries := 50
		c.Logger.Info(conn.State().String())
		for conn.State() != seqs.StateEstablished && retries > 0 {
			time.Sleep(100 * time.Millisecond)
			retries--
		}
		if retries == 0 {
			c.Logger.Info("socket:no-establish")
			closeConn("did not establish connection")
			continue
		}

		// We start MQTT connect with a deadline on the socket.
		c.Logger.Info("mqtt:start-connecting")
		lcd.Send(lcdMessages, "MQTT Connect", "Authenticating")
		conn.SetDeadline(time.Now().Add(c.Timeout))
		err = mqttClient.StartConnect(conn, &varconn)
		if err != nil {
			c.Logger.Error("mqtt:start-connect-failed", slog.String("reason", err.Error()))
			lcd.Send(lcdMessages, "Connect Failed", err.Error()[:min(len(err.Error()), 16)])
			closeConn("connect failed")
			continue
		}
		retries = 50
		for retries > 0 && !mqttClient.IsConnected() {
			time.Sleep(100 * time.Millisecond)
			err = mqttClient.HandleNext()
			if err != nil {
				c.Logger.Error("mqtt:handle-next-failed", slog.String("err", err.Error()))
			}
			retries--
		}
		if !mqttClient.IsConnected() {
			c.Logger.Error("mqtt:connect-failed", slog.Any("reason", mqttClient.Err()))
			lcd.Send(lcdMessages, "Connect Failed", "Timed out")
			closeConn("connect timed out")
			continue
		}

		lcd.Send(lcdMessages, "MQTT Connected", "Publishing...")

		heartbeat := time.NewTicker(c.HeartbeatInterval)
		defer heartbeat.Stop()
		for mqttClient.IsConnected() {
			select {
			case reading := <-readings:
				payload, err := json.Marshal(reading)
				if err != nil {
					c.Logger.Error("mqtt:marshal-failed", slog.Any("reason", err))
					continue
				}
				conn.SetDeadline(time.Now().Add(c.Timeout))
				pubVar.PacketIdentifier = uint16(rng.Uint32())
				err = mqttClient.PublishPayload(pubFlags, pubVar, payload)
				if err != nil {
					c.Logger.Error("mqtt:publish-failed", slog.Any("reason", err))
					continue
				}
				c.Logger.Info("published message",
					slog.Uint64("packetID", uint64(pubVar.PacketIdentifier)),
				)
				err = mqttClient.HandleNext()
				if err != nil {
					c.Logger.Error("mqtt:handle-next-failed", slog.String("err", err.Error()))
					continue
				}
			case <-heartbeat.C:
				// If we haven't read any sensor readings from the channel since the last heartbeat interval,
				// ping the MQTT broken to keep the connnection alive.
				err = mqttClient.HandleNext()
				if err != nil {
					c.Logger.Error("mqtt:handle-next-failed", slog.String("err", err.Error()))
					continue
				}
			default:
				// If we've got nothing to do, release the thread so other go routines can run.
				// We only need to do this because TinyGo runs on a single core
				// https://tinygo.org/docs/guides/tips-n-tricks/
				runtime.Gosched()
			}

		}

		c.Logger.Error("mqtt:disconnected", slog.Any("reason", mqttClient.Err()))
		lcd.Send(lcdMessages, "Disconnected", "Reconnecting...")
		closeConn("disconnected")
		runtime.Gosched()
	}
}

// splitHostPort splits a host:port string into separate host and port components.
// Returns an error if the format is invalid.
func splitHostPort(addr string) (host, port string, err error) {
	// Find the last colon to support IPv6 addresses
	colonIdx := -1
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			colonIdx = i
			break
		}
	}

	if colonIdx == -1 {
		return "", "", errors.New("missing port in address")
	}

	host = addr[:colonIdx]
	port = addr[colonIdx+1:]

	if host == "" {
		return "", "", errors.New("empty host")
	}
	if port == "" {
		return "", "", errors.New("empty port")
	}

	return host, port, nil
}

// parsePort converts a port string to uint16.
// Returns 0 if parsing fails (caller should validate).
func parsePort(portStr string) uint16 {
	var port uint16
	for i := 0; i < len(portStr); i++ {
		if portStr[i] < '0' || portStr[i] > '9' {
			return 0
		}
		port = port*10 + uint16(portStr[i]-'0')
	}
	return port
}
