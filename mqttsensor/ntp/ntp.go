package ntp

import (
	"errors"
	"log/slog"
	"runtime"
	"time"

	"github.com/harveysanders/picoplayground/mqttsensor/cyw43439"
	"github.com/soypat/seqs/eth/ntp"
	"github.com/soypat/seqs/stacks"
)

// SyncTime synchronizes the system time using NTP.
// It resolves an NTP server via DNS, performs an NTP time sync,
// and adjusts the system clock accordingly.
//
// Returns nil on success, error if sync fails.
func SyncTime(
	stack *stacks.PortStack,
	resolver *cyw43439.Resolver,
	routerHWAddr [6]byte,
	logger *slog.Logger,
) error {
	// Resolve NTP server address
	logger.Info("ntp:resolving pool.ntp.org")
	addrs, err := resolver.LookupNetIP("pool.ntp.org")
	if err != nil {
		return errors.New("ntp dns lookup:" + err.Error())
	}
	if len(addrs) == 0 {
		return errors.New("ntp dns lookup: no addresses returned")
	}

	ntpAddr := addrs[0]
	logger.Info("ntp:resolved", slog.String("addr", ntpAddr.String()))

	// Create NTP client
	ntpc := stacks.NewNTPClient(stack, ntp.ClientPort)
	err = ntpc.BeginDefaultRequest(routerHWAddr, ntpAddr)
	if err != nil {
		return errors.New("ntp request:" + err.Error())
	}

	// Wait for NTP response (with timeout)
	logger.Info("ntp:waiting for response")
	const timeout = 5 * time.Second
	const pollInterval = 100 * time.Millisecond
	start := time.Now()
	for !ntpc.IsDone() {
		if time.Since(start) > timeout {
			return errors.New("ntp timeout")
		}
		time.Sleep(pollInterval)
	}

	// Calculate and apply time offset
	// NTP base time is 1900-01-01T00:00:00.000Z
	offset := ntp.BaseTime().
		// Add the diff from actual now (from NTP) and 1900-01-01
		// This offset is going to stop working around 2036-02-07
		// https://docs.ntpsec.org/latest/rollover.html
		Add(ntpc.Offset()).
		// Subtract the device's time (some time close to 0 unix time (1970-01-01...))
		Sub(time.Now())

		// Use the offset to set the device's time
	runtime.AdjustTimeOffset(int64(offset))
	return nil
}
