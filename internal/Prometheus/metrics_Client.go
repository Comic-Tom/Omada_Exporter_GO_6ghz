package Prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"omada_exporter_go/internal/Omada/Model/Client"
)

// ── Label names ────────────────────────────────────────────────────────────────

const (
	label_clientMAC        string = "mac"
	label_clientName       string = "clientName"
	label_clientIP         string = "ip"
	label_clientVendor     string = "vendor"
	label_clientDeviceType string = "clientDeviceType"
	label_clientNetwork    string = "networkName"
	label_clientWireless   string = "wireless"
	label_clientVLAN       string = "vlan"

	// Wireless labels
	label_clientSSID     string = "ssid"
	label_clientBand     string = "band"
	label_clientApMAC    string = "apMac"
	label_clientApName   string = "apName"
	label_clientWifiStd  string = "wifiStandard"
	label_clientMLO      string = "mlo"

	// Wired labels
	label_clientSwitchMAC  string = "switchMac"
	label_clientSwitchName string = "switchName"
	label_clientSwitchPort string = "switchPort"
)

// ── Label sets ─────────────────────────────────────────────────────────────────

// Base labels present on ALL client metrics (wired + wireless)
var clientBaseLabels = []string{
	label_clientMAC,
	label_clientName,
	label_clientIP,
	label_clientVendor,
	label_clientDeviceType,
	label_clientNetwork,
	label_clientWireless,
	label_clientVLAN,
}

// Extra labels for wireless-specific metrics
var clientWirelessLabels = append(clientBaseLabels,
	label_clientSSID,
	label_clientBand,
	label_clientApMAC,
	label_clientApName,
	label_clientWifiStd,
	label_clientMLO,
)

// Extra labels for wired-specific metrics
var clientWiredLabels = append(clientBaseLabels,
	label_clientSwitchMAC,
	label_clientSwitchName,
	label_clientSwitchPort,
)

// ── Metrics ────────────────────────────────────────────────────────────────────

var (
	// All clients
	client_info = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_info",
			Help: "Information about a connected client (value always 1)",
		},
		clientBaseLabels,
	)
	client_traffic_down_bytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_traffic_down_bytes_total",
			Help: "Total bytes downloaded by the client in the current session",
		},
		clientBaseLabels,
	)
	client_traffic_up_bytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_traffic_up_bytes_total",
			Help: "Total bytes uploaded by the client in the current session",
		},
		clientBaseLabels,
	)
	client_uptime_seconds = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_uptime_seconds",
			Help: "How long the client has been connected in seconds",
		},
		clientBaseLabels,
	)

	// Wireless clients
	client_signal_level = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_signal_level_dbm",
			Help: "Wireless signal level (RSSI) of the client in dBm",
		},
		clientWirelessLabels,
	)
	client_snr = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_snr_db",
			Help: "Signal-to-noise ratio of the client's wireless connection in dB",
		},
		clientWirelessLabels,
	)
	client_rx_rate_bps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_rx_rate_bps",
			Help: "Current RX (download) link rate of the client in bits per second",
		},
		clientWirelessLabels,
	)
	client_tx_rate_bps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_tx_rate_bps",
			Help: "Current TX (upload) link rate of the client in bits per second",
		},
		clientWirelessLabels,
	)

	// MLO per-link metrics (one series per client per active radio link)
	client_mlo_link_rx_rate_bps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_mlo_link_rx_rate_bps",
			Help: "Per-link RX rate for MLO clients in bits per second",
		},
		[]string{label_clientMAC, label_clientName, label_clientApMAC, label_clientApName, "band", "wifiStandard"},
	)
	client_mlo_link_tx_rate_bps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_mlo_link_tx_rate_bps",
			Help: "Per-link TX rate for MLO clients in bits per second",
		},
		[]string{label_clientMAC, label_clientName, label_clientApMAC, label_clientApName, "band", "wifiStandard"},
	)
	client_mlo_link_snr = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_mlo_link_snr_db",
			Help: "Per-link SNR for MLO clients in dB",
		},
		[]string{label_clientMAC, label_clientName, label_clientApMAC, label_clientApName, "band", "wifiStandard"},
	)
	client_mlo_link_signal_level = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_mlo_link_signal_level_dbm",
			Help: "Per-link RSSI for MLO clients in dBm",
		},
		[]string{label_clientMAC, label_clientName, label_clientApMAC, label_clientApName, "band", "wifiStandard"},
	)

	// Wired clients
	client_wired_link_speed_bps = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_wired_link_speed_bps",
			Help: "Negotiated wired link speed of the client in bits per second (derived from switch port speed)",
		},
		clientWiredLabels,
	)
)

// ── Helpers ────────────────────────────────────────────────────────────────────

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func getClientBaseLabels(c Client.Client) prometheus.Labels {
	return prometheus.Labels{
		label_clientMAC:        c.MAC,
		label_clientName:       c.DisplayName(),
		label_clientIP:         c.IP,
		label_clientVendor:     c.Vendor,
		label_clientDeviceType: c.DeviceType,
		label_clientNetwork:    c.NetworkName,
		label_clientWireless:   boolToString(c.Wireless),
		label_clientVLAN:       fmt.Sprintf("%d", c.VLAN),
	}
}

func getClientWirelessLabels(c Client.Client) prometheus.Labels {
	base := getClientBaseLabels(c)
	base[label_clientSSID] = c.SSID
	base[label_clientBand] = c.RadioBand()
	base[label_clientApMAC] = c.ApMAC
	base[label_clientApName] = c.ApName
	base[label_clientWifiStd] = c.WifiMode.String()
	base[label_clientMLO] = boolToString(c.IsMLO())
	return base
}

func getClientWiredLabels(c Client.Client) prometheus.Labels {
	base := getClientBaseLabels(c)
	base[label_clientSwitchMAC] = c.SwitchMAC
	base[label_clientSwitchName] = c.SwitchName
	base[label_clientSwitchPort] = c.SwitchPortID()
	return base
}

// ── Main expose function ───────────────────────────────────────────────────────

// ExposeClientMetrics resets and republishes all per-client Prometheus metrics.
func ExposeClientMetrics(clients []Client.Client) {
	// Reset everything each scrape so disconnected clients disappear
	client_info.Reset()
	client_traffic_down_bytes.Reset()
	client_traffic_up_bytes.Reset()
	client_uptime_seconds.Reset()
	client_signal_level.Reset()
	client_snr.Reset()
	client_rx_rate_bps.Reset()
	client_tx_rate_bps.Reset()
	client_mlo_link_rx_rate_bps.Reset()
	client_mlo_link_tx_rate_bps.Reset()
	client_mlo_link_snr.Reset()
	client_mlo_link_signal_level.Reset()
	client_wired_link_speed_bps.Reset()

	for _, c := range clients {
		base := getClientBaseLabels(c)

		client_info.With(base).Set(1)
		client_traffic_down_bytes.With(base).Set(c.TrafficDown)
		client_traffic_up_bytes.With(base).Set(c.TrafficUp)
		client_uptime_seconds.With(base).Set(c.Uptime)

		if c.Wireless {
			wl := getClientWirelessLabels(c)

			client_signal_level.With(wl).Set(float64(c.SignalLevel))
			client_snr.With(wl).Set(float64(c.SNR))
			client_rx_rate_bps.With(wl).Set(c.RxRate * 1000) // kbps → bps
			client_tx_rate_bps.With(wl).Set(c.TxRate * 1000)

			// MLO per-link breakdown
			if c.IsMLO() {
				for _, link := range c.MultiLink {
					mloLabels := prometheus.Labels{
						label_clientMAC:    c.MAC,
						label_clientName:   c.DisplayName(),
						label_clientApMAC:  c.ApMAC,
						label_clientApName: c.ApName,
						"band":             link.RadioBand(),
						"wifiStandard":     link.WifiMode.String(),
					}
					client_mlo_link_rx_rate_bps.With(mloLabels).Set(link.RxRate * 1000)
					client_mlo_link_tx_rate_bps.With(mloLabels).Set(link.TxRate * 1000)
					client_mlo_link_snr.With(mloLabels).Set(float64(link.SNR))
					client_mlo_link_signal_level.With(mloLabels).Set(float64(link.RSSI))
				}
			}
		} else {
			// Wired client — link speed comes from switch port
			// We set 0 here as a placeholder; actual speed requires a switch port lookup.
			// The switchPort label still gives useful topology info in Grafana.
			wired := getClientWiredLabels(c)
			client_wired_link_speed_bps.With(wired).Set(0)
		}
	}
}
