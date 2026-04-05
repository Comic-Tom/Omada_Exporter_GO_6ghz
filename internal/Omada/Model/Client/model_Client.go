package Client

import (
	"fmt"
	"omada_exporter_go/internal/Omada/Enum"
)

const path_OpenApiClients = "/openapi/v1/{omadaID}/sites/{siteID}/clients"

// MultiLinkEntry represents a single radio link in an MLO connection.
// Note: the official multiLink wifiMode docs only define 0-7; values 8/9
// are observed in practice for 802.11be (WiFi 7) links.
type MultiLinkEntry struct {
	RadioID     int           `json:"radioId"`
	WifiMode    Enum.WifiMode `json:"wifiMode"`
	Channel     int           `json:"channel"`
	RxRate      float64       `json:"rxRate"`      // kbps
	TxRate      float64       `json:"txRate"`      // kbps
	SNR         int           `json:"snr"`
	RSSI        int           `json:"rssi"`        // dBm
	SignalLevel int           `json:"signalLevel"` // 0-100 percentage
}

func (m MultiLinkEntry) RadioBand() string {
	return RadioBandFromID(m.RadioID)
}

// WifiStandardString returns a band-aware WiFi standard label for this link.
// For wifiMode 9 (WiFi 7), the band cannot be determined from wifiMode alone —
// radioId is required to distinguish 5GHz (radioId=1/2) from 6GHz (radioId=3).
func (m MultiLinkEntry) WifiStandardString() string {
	return m.WifiMode.StringWithBand(m.RadioID)
}

// RadioBandFromID maps Omada's radioId to a band string.
// 0=2.4GHz, 1=5GHz-1, 2=5GHz-2, 3=6GHz
func RadioBandFromID(radioID int) string {
	switch radioID {
	case 0:
		return "2.4 GHz"
	case 1:
		return "5.0 GHz"
	case 2:
		return "5.0 GHz-2"
	case 3:
		return "6.0 GHz"
	default:
		return "unknown"
	}
}

// Client represents a single network client.
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

	// Connection device type: "ap", "switch", "gateway"
	ConnectDevType string `json:"connectDevType"`

	// Traffic
	TrafficDown float64 `json:"trafficDown"` // bytes
	TrafficUp   float64 `json:"trafficUp"`   // bytes

	// Wired — switch-connected
	SwitchMAC  string `json:"switchMac"`
	SwitchName string `json:"switchName"`
	SwitchPort int    `json:"port"`

	// Wired — gateway-connected
	GatewayMAC  string `json:"gatewayMac"`
	GatewayName string `json:"gatewayName"`

	// Wireless top-level fields (use PrimaryLink() for signal/rate data)
	SSID     string        `json:"ssid"`
	ApMAC    string        `json:"apMac"`
	ApName   string        `json:"apName"`
	RadioID  int           `json:"radioId"`
	WifiMode Enum.WifiMode `json:"wifiMode"`
	RxRate   float64       `json:"rxRate"` // kbps
	TxRate   float64       `json:"txRate"` // kbps
	RSSI     int           `json:"rssi"`   // dBm
	SNR      int           `json:"snr"`

	// MLO links
	MultiLink []MultiLinkEntry `json:"multiLink"`
}

// DisplayName returns the best human-readable name.
func (c Client) DisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.HostName != "" {
		return c.HostName
	}
	return c.MAC
}

// IsWifi7 returns true if any link uses 802.11be (wifiMode 8 or 9).
// wifiMode 8 = WiFi 7 2.4GHz, wifiMode 9 = WiFi 7 5GHz/6GHz.
func (c Client) IsWifi7() bool {
	if c.WifiMode == Enum.WifiMode_11beg || c.WifiMode == Enum.WifiMode_11bea {
		return true
	}
	for _, link := range c.MultiLink {
		if link.WifiMode == Enum.WifiMode_11beg || link.WifiMode == Enum.WifiMode_11bea {
			return true
		}
	}
	return false
}

// PrimaryLink returns the best available link for signal/rate metrics.
// For MLO clients, returns the link with the highest RSSI.
// For non-MLO wireless clients, synthesises from top-level fields.
func (c Client) PrimaryLink() *MultiLinkEntry {
	if len(c.MultiLink) == 0 {
		if !c.Wireless {
			return nil
		}
		return &MultiLinkEntry{
			RadioID:  c.RadioID,
			WifiMode: c.WifiMode,
			RxRate:   c.RxRate,
			TxRate:   c.TxRate,
			RSSI:     c.RSSI,
			SNR:      c.SNR,
		}
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
	return "unknown"
}

// WifiStandard returns a band-aware WiFi standard label for the primary link.
// For wifiMode 9 (WiFi 7 5GHz/6GHz), radioId is used to distinguish the band.
func (c Client) WifiStandard() string {
	if link := c.PrimaryLink(); link != nil {
		return link.WifiMode.StringWithBand(link.RadioID)
	}
	return c.WifiMode.StringWithBand(c.RadioID)
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

// LinkRxRate returns the RX rate in kbps of the primary link.
func (c Client) LinkRxRate() float64 {
	if link := c.PrimaryLink(); link != nil {
		return link.RxRate
	}
	return c.RxRate
}

// LinkTxRate returns the TX rate in kbps of the primary link.
func (c Client) LinkTxRate() float64 {
	if link := c.PrimaryLink(); link != nil {
		return link.TxRate
	}
	return c.TxRate
}

// IsMLO returns true if this client has multiple active radio links.
func (c Client) IsMLO() bool {
	active := 0
	for _, l := range c.MultiLink {
		if l.RxRate > 0 || l.TxRate > 0 || l.SNR > 0 {
			active++
		}
	}
	return active > 1
}

// SwitchPortID returns the switch port as a string.
func (c Client) SwitchPortID() string {
	if c.SwitchPort == 0 {
		return Enum.NotApplicable_String
	}
	return fmt.Sprintf("%d", c.SwitchPort)
}

// ConnectDevice returns the name of the device this client is connected to.
func (c Client) ConnectDevice() string {
	switch c.ConnectDevType {
	case "switch":
		return c.SwitchName
	case "gateway":
		return c.GatewayName
	case "ap":
		return c.ApName
	default:
		return ""
	}
}

// ConnectDeviceMAC returns the MAC of the device this client is connected to.
func (c Client) ConnectDeviceMAC() string {
	switch c.ConnectDevType {
	case "switch":
		return c.SwitchMAC
	case "gateway":
		return c.GatewayMAC
	case "ap":
		return c.ApMAC
	default:
		return ""
	}
}
