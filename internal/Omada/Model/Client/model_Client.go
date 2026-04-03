package Client

import (
	"fmt"
	"omada_exporter_go/internal/Omada/Enum"
)

const path_OpenApiClients = "/openapi/v1/{omadaID}/sites/{siteID}/clients"

// MultiLinkEntry represents a single radio link in an MLO (Multi-Link Operation) connection.
// WiFi 7 clients may have multiple entries — one per active radio band.
type MultiLinkEntry struct {
	RadioID     int           `json:"radioId"`     // 0=2.4GHz, 1=5GHz, 3=6GHz (Omada skips 2)
	WifiMode    Enum.WifiMode `json:"wifiMode"`
	Channel     int           `json:"channel"`
	RxRate      float64       `json:"rxRate"`      // kbps
	TxRate      float64       `json:"txRate"`      // kbps
	SNR         int           `json:"snr"`
	RSSI        int           `json:"rssi"`        // actual dBm signal strength
	SignalLevel int           `json:"signalLevel"` // 0-100 percentage (not dBm)
}

// RadioBand maps Omada's radioId to a human-readable band string.
// Note: Omada uses 0=2.4GHz, 1=5GHz, 3=6GHz — there is no radioId 2.
func RadioBandFromID(radioID int) string {
	switch radioID {
	case 0:
		return "2.4 GHz"
	case 1:
		return "5.0 GHz"
	case 3:
		return "6.0 GHz"
	default:
		return "unknown"
	}
}

func (m MultiLinkEntry) RadioBand() string {
	return RadioBandFromID(m.RadioID)
}

// Client represents a single network client connected to the Omada controller.
type Client struct {
	// Identity
	MAC         string `json:"mac"`
	Name        string `json:"name"`
	HostName    string `json:"hostName"`
	IP          string `json:"ip"`
	Vendor      string `json:"vendor"`
	DeviceType  string `json:"deviceType"`
	NetworkName string `json:"networkName"`
	VLAN        int    `json:"vid"`

	// Connection state
	Active   bool    `json:"active"`
	Wireless bool    `json:"wireless"`
	Uptime   float64 `json:"uptime"` // seconds

	// Traffic (session totals, in bytes)
	TrafficDown float64 `json:"trafficDown"`
	TrafficUp   float64 `json:"trafficUp"`

	// Wired-specific
	SwitchMAC  string `json:"switchMac"`
	SwitchName string `json:"switchName"`
	SwitchPort int    `json:"port"`

	// Wireless top-level fields (unreliable for MLO — use MultiLink instead)
	SSID     string        `json:"ssid"`
	ApMAC    string        `json:"apMac"`
	ApName   string        `json:"apName"`
	RadioID  int           `json:"radioId"`
	WifiMode Enum.WifiMode `json:"wifiMode"`

	// Top-level rates (kbps) — use PrimaryLink() for accurate per-band data
	RxRate float64 `json:"rxRate"`
	TxRate float64 `json:"txRate"`

	// WiFi 7 MLO — one entry per active radio link
	MultiLink []MultiLinkEntry `json:"multiLink"`
}

// DisplayName returns the best available human-readable name for the client.
func (c Client) DisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.HostName != "" {
		return c.HostName
	}
	return c.MAC
}

// PrimaryLink returns the best available single link for signal/rate metrics.
// For MLO clients, returns the link with the highest RSSI (best signal).
// For non-MLO clients, returns the single multiLink entry if present,
// falling back to synthesising from top-level fields.
func (c Client) PrimaryLink() *MultiLinkEntry {
	if len(c.MultiLink) == 0 {
		return nil
	}
	best := &c.MultiLink[0]
	for i := range c.MultiLink {
		if c.MultiLink[i].RSSI > best.RSSI {
			best = &c.MultiLink[i]
		}
	}
	return best
}

// RadioBand returns the band of the primary link.
func (c Client) RadioBand() string {
	if link := c.PrimaryLink(); link != nil {
		return link.RadioBand()
	}
	return RadioBandFromID(c.RadioID)
}

// WifiStandard returns the WiFi standard of the primary link.
func (c Client) WifiStandard() string {
	if link := c.PrimaryLink(); link != nil {
		return link.WifiMode.String()
	}
	return c.WifiMode.String()
}

// SignalRSSI returns the RSSI (dBm) of the primary link.
func (c Client) SignalRSSI() float64 {
	if link := c.PrimaryLink(); link != nil {
		return float64(link.RSSI)
	}
	return 0
}

// SignalSNR returns the SNR (dB) of the primary link.
func (c Client) SignalSNR() float64 {
	if link := c.PrimaryLink(); link != nil {
		return float64(link.SNR)
	}
	return 0
}

// LinkRxRate returns the RX rate (kbps) of the primary link.
func (c Client) LinkRxRate() float64 {
	if link := c.PrimaryLink(); link != nil {
		return link.RxRate
	}
	return c.RxRate
}

// LinkTxRate returns the TX rate (kbps) of the primary link.
func (c Client) LinkTxRate() float64 {
	if link := c.PrimaryLink(); link != nil {
		return link.TxRate
	}
	return c.TxRate
}

// IsMLO returns true if this client is connected via MLO on multiple active links.
func (c Client) IsMLO() bool {
	active := 0
	for _, l := range c.MultiLink {
		if l.RxRate > 0 || l.TxRate > 0 || l.SNR > 0 {
			active++
		}
	}
	return active > 1
}

// SwitchPortID returns the switch port as a string label.
func (c Client) SwitchPortID() string {
	if c.SwitchPort == 0 {
		return Enum.NotApplicable_String
	}
	return fmt.Sprintf("%d", c.SwitchPort)
}
