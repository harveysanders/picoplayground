package mqtt

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math/rand"
	"net/netip"
	"runtime"
	"time"

	"github.com/harveysanders/picoplayground/mqttsensor/cyw43439"
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
	Voltage   float32
	RawValue  uint16
	SinceBoot time.Duration
}

type Client struct {
	ID                string
	Timeout           time.Duration
	TCPBufSize        int
	Logger            *slog.Logger
	HeartbeatInterval time.Duration
}

func (c Client) ConnectAndPublish(addr string, readings <-chan SensorReading) error {
	_, stack, _, err := cyw43439.SetupWithDHCP(cyw43439.SetupConfig{
		Hostname: c.ID,
		Logger:   c.Logger,
		TCPPorts: 1, // For HTTP over TCP.
		UDPPorts: 1, // For DNS.
	})

	start := time.Now()
	if err != nil {
		return errors.New("setup DHCP:" + err.Error())
	}
	svAddr, err := netip.ParseAddrPort(addr)
	if err != nil {
		return errors.New("parsing server address:" + err.Error())
	}
	// Resolver router's hardware address to dial outside our network to internet.
	serverHWAddr, err := cyw43439.ResolveHardwareAddr(stack, svAddr.Addr())
	if err != nil {
		return errors.New("router hwaddr resolving:" + err.Error())
	}
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
	client := mqtt.NewClient(cfg)

	// Connection loop for TCP+MQTT.
	for {
		random := rng.Uint32()
		c.Logger.Info("socket:listen")
		err = conn.OpenDialTCP(clientAddr.Port(), serverHWAddr, svAddr, seqs.Value(random))
		if err != nil {
			panic("socket dial:" + err.Error())
		}
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
		conn.SetDeadline(time.Now().Add(c.Timeout))
		err = client.StartConnect(conn, &varconn)
		if err != nil {
			c.Logger.Error("mqtt:start-connect-failed", slog.String("reason", err.Error()))
			closeConn("connect failed")
			continue
		}
		retries = 50
		for retries > 0 && !client.IsConnected() {
			time.Sleep(100 * time.Millisecond)
			err = client.HandleNext()
			if err != nil {
				c.Logger.Error("mqtt:handle-next-failed", slog.String("err", err.Error()))
			}
			retries--
		}
		if !client.IsConnected() {
			c.Logger.Error("mqtt:connect-failed", slog.Any("reason", client.Err()))
			closeConn("connect timed out")
			continue
		}

		for client.IsConnected() {
			select {
			case reading := <-readings:
				payload, err := json.Marshal(reading)
				if err != nil {
					c.Logger.Error("mqtt:marshal-failed", slog.Any("reason", err))
					continue
				}
				conn.SetDeadline(time.Now().Add(c.Timeout))
				pubVar.PacketIdentifier = uint16(rng.Uint32())
				err = client.PublishPayload(pubFlags, pubVar, payload)
				if err != nil {
					c.Logger.Error("mqtt:publish-failed", slog.Any("reason", err))
					continue
				}
				c.Logger.Info("published message", slog.Uint64("packetID", uint64(pubVar.PacketIdentifier)))
				err = client.HandleNext()
				if err != nil {
					c.Logger.Error("mqtt:handle-next-failed", slog.String("err", err.Error()))
					continue
				}
			case <-time.After(c.HeartbeatInterval):
				// If we haven't read any sensor readings from the channel since the last heartbeat interval,
				// ping the MQTT broken to keep the connnection alive.
				err = client.HandleNext()
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

		c.Logger.Error("mqtt:disconnected", slog.Any("reason", client.Err()))
		closeConn("disconnected")
		runtime.Gosched()
	}
}
