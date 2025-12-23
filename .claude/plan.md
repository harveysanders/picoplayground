# NTP Time Synchronization Implementation Plan

## Overview
Add NTP time synchronization to the Pico2 sensor project to get accurate timestamps instead of just "time since boot". The implementation will sync time once at startup after WiFi connects.

## Dependencies
All required packages are already in go.mod:
- `github.com/soypat/seqs/eth/ntp` - NTP protocol implementation
- `github.com/soypat/seqs/stacks` - Network stack with NTP client
- `runtime` (stdlib) - For `AdjustTimeOffset()` to set system time

## Implementation Steps

### 1. Create NTP sync function in mqtt package
**File**: `mqttsensor/mqtt/ntp.go` (new file)

Add a `SyncTime()` function that:
- Takes the network stack, resolver, and router hardware address (already available in ConnectAndPublish)
- Resolves "pool.ntp.org" to an IP address using DNS
- Creates NTP client using `stacks.NewNTPClient(stack, ntp.ClientPort)`
- Sends NTP request with `ntpc.BeginDefaultRequest(routerhw, ntpaddr)`
- Waits for completion (with timeout)
- Calculates time offset and adjusts system time using `runtime.AdjustTimeOffset()`
- **Returns bool (success/failure)** - true if time was synced, false if NTP failed

### 2. Integrate NTP sync into ConnectAndPublish
**File**: `mqttsensor/mqtt/client.go`

In the `ConnectAndPublish()` function:
- After successful DHCP setup (line ~53)
- After creating the DNS resolver (line ~67)
- Call `SyncTime()` before entering the MQTT connection loop
- Send LCD message showing time sync status
- **If successful, set `c.TimeSyncedAt = time.Now()`**
- **If failed, leave `c.TimeSyncedAt` as zero time (default)**
- Log success/failure but don't block if NTP fails

### 3. Update SensorReading struct
**File**: `mqttsensor/mqtt/client.go`

Add a `Timestamp` field to the `SensorReading` struct:
```go
type SensorReading struct {
    Voltage   float32
    RawValue  uint16
    SinceBoot time.Duration
    Timestamp time.Time  // Go zero time if NTP sync failed, actual time if succeeded
}
```

### 4. Update Client struct to track NTP sync time
**File**: `mqttsensor/mqtt/client.go`

Add a field to track when time was synced:
```go
type Client struct {
    ID                string
    Timeout           time.Duration
    TCPBufSize        int
    Logger            *slog.Logger
    HeartbeatInterval time.Duration
    TimeSyncedAt      time.Time  // When NTP sync occurred. Zero if never synced.
}
```

Benefits:
- `TimeSyncedAt.IsZero()` tells us if NTP sync succeeded
- `time.Since(c.TimeSyncedAt)` enables periodic re-sync later if needed

### 5. Update sensor reading creation
**File**: `mqttsensor/main.go`

When creating sensor readings (line ~124), conditionally set timestamp:
```go
reading := mqtt.SensorReading{
    Voltage:   voltage,
    RawValue:  val,
    SinceBoot: time.Since(start),
}

// Only set Timestamp if NTP sync succeeded
// Check if c.TimeSyncedAt is not zero (meaning NTP sync happened)
if !c.TimeSyncedAt.IsZero() {
    reading.Timestamp = time.Now()
}

sensorReadings <- reading
```

Main will access `c.TimeSyncedAt` to determine if time is synced. Since `TimeSyncedAt` is written once at startup and read repeatedly, and `time.Time` operations are safe for concurrent read, this approach is simple and effective.

## Design Decisions

### NTP Server
- Use "pool.ntp.org" (public NTP pool)
- Resolve via DNS at runtime (not hardcoded IP)

### Error Handling
- NTP sync failure is non-fatal
- Log error and continue with operation
- If sync fails, `Timestamp` field will be Go zero time (`time.Time{}`), not Unix epoch
- Consumers can use `timestamp.IsZero()` to check if time is valid

### Zero Time vs Unix Epoch
- **Important**: On embedded systems, if we don't sync time, `time.Now()` returns Unix epoch (Jan 1, 1970), NOT Go zero time
- We explicitly set `Timestamp` to zero time (`time.Time{}`) when NTP fails
- This makes it easy to detect: `if !timestamp.IsZero() { /* use timestamp */ }`

### Sync Frequency
- One-time sync at boot
- No periodic re-sync (user's tolerance of "few hundred ms" makes this acceptable)
- If needed later, can add hourly re-sync

### LCD Feedback
- Show "Syncing time..." message during NTP request
- Show "Time synced" on success or continue to next status on failure

## Files to Create
- `mqttsensor/mqtt/ntp.go` - New file with SyncTime() function

## Files to Modify
- `mqttsensor/mqtt/client.go`:
  - Add `Timestamp` field to `SensorReading`
  - Add `TimeSyncedAt` field to `Client` struct
  - Integrate NTP sync call in `ConnectAndPublish()`
- `mqttsensor/main.go`:
  - Conditionally set `Timestamp` based on `c.TimeSyncedAt`

## Testing Considerations
- NTP sync should complete within 2-3 seconds
- If NTP server is unreachable, should timeout and continue
- Verify `time.Now()` returns correct time after sync
- Check that sensor readings contain proper timestamps when synced
- Verify `Timestamp.IsZero()` returns true when NTP fails
