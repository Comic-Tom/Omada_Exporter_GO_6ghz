package Client

import (
	"fmt"
	"omada_exporter_go/internal/Omada/Enum"
)

const path_OpenApiClients = "/openapi/v1/{omadaID}/sites/{siteID}/clients"

// MultiLinkEntry represents a single radio link in an MLO (Multi-Link Operation) connection.
// WiFi 7 clients may have multiple entries — one per active radio band.
type MultiLinkEntry struct {
	RadioID  int           `json:"radioId"`  // 0=2.4GHz, 1=5GHz, 2=6GHz
	WifiMode Enum.WifiMode `json:"wifiMode"` // WiFi standard
	Channel  int           `json:"channel"`
	RxRate   float64       `json:"rxRate"` // kbps
	TxRate   float64       `json:"txRate"` // kbps
	SNR      int           `json:"snr"`
	RSSI     int           `json:"rssi"` // signal in dBm
}

// RadioBand returns a human-readable band for a MultiLinkEntry.
// MLO radioId mapping: 0=2.4GHz, 1=5GHz, 2=5GHz-2 (some APs), 3=6GHz
func (m MultiLinkEntry) RadioBand() string {
	switch m.RadioID {
	case 0:
		return "2.4 GHz"
	case 1:
		return "5.0 GHz"
	case 2:
		return "5.0 GHz (2)"
	case 3:
		return "6.0 GHz"
	default:
		return "unknown"
	}
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
	SwitchPort int    `json:"port"` // physical port number on the switch

	// Wireless-specific (single-link / non-MLO)
	SSID        string        `json:"ssid"`
	ApMAC       string        `json:"apMac"`
	ApName      string        `json:"apName"`
	RadioID     int           `json:"radioId"`     // 0=2.4GHz, 1=5GHz, 2=6GHz
	WifiMode    Enum.WifiMode `json:"wifiMode"`    // WiFi standard
	SignalLevel int           `json:"signalLevel"` // dBm (RSSI)
	SNR         int           `json:"snr"`
	RxRate      float64       `json:"rxRate"` // kbps
	TxRate      float64       `json:"txRate"` // kbps

	// WiFi 7 MLO — multiple radio links active simultaneously
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

// RadioBand returns a human-readable band string for a single-link wireless client.
// radioId mapping: 0=2.4GHz, 1=5GHz, 2=5GHz-2 (some APs), 3=6GHz
func (c Client) RadioBand() string {
	switch c.RadioID {
	case 0:
		return "2.4 GHz"
	case 1:
		return "5.0 GHz"
	case 2:
		return "5.0 GHz (2)"
	case 3:
		return "6.0 GHz"
	default:
		return "unknown"
	}
}

// SwitchPortID returns the switch port as a string label.
func (c Client) SwitchPortID() string {
	if c.SwitchPort == 0 {
		return Enum.NotApplicable_String
	}
	return fmt.Sprintf("%d", c.SwitchPort)
}

// IsMLO returns true if this client is connected via WiFi 7 MLO on multiple links.
func (c Client) IsMLO() bool {
	return len(c.MultiLink) > 1
}
