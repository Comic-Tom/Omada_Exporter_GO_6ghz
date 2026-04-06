package Prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"omada_exporter_go/internal/Omada/Model/Client"
)

const (
	label_clientMAC         string = "mac"
	label_clientName        string = "clientName"
	label_clientIP          string = "ip"
	label_clientVendor      string = "vendor"
	label_clientDeviceType  string = "clientDeviceType"
	label_clientNetwork     string = "networkName"
	label_clientWireless    string = "wireless"
	label_clientVLAN        string = "vlan"
	label_clientSSID        string = "ssid"
	label_clientBand        string = "band"
	label_clientApMAC       string = "apMac"
	label_clientApName      string = "apName"
	label_clientWifiStd     string = "wifiStandard"
	label_clientMLO         string = "mlo"
	label_clientConnectDev  string = "connectDevice"
	label_clientConnectMAC  string = "connectDevMac"
	label_clientConnectType string = "connectDevType"
	label_clientSwitchPort  string = "switchPort"
)

var clientBaseLabels = []string{
	label_clientMAC, label_clientName, label_clientIP,
	label_clientVendor, label_clientDeviceType,
	label_clientNetwork, label_clientWireless, label_clientVLAN,
}

var clientWirelessLabels = append([]string{}, append(clientBaseLabels,
	label_clientSSID, label_clientBand, label_clientApMAC,
	label_clientApName, label_clientWifiStd, label_clientMLO,
)...)

// Wired labels include connect device info (switch or gateway) and port.
var clientWiredLabels = append([]string{}, append(clientBaseLabels,
	label_clientConnectDev, label_clientConnectMAC,
	label_clientConnectType, label_clientSwitchPort,
)...)

var mloLinkLabels = []string{
	label_clientMAC, label_clientName,
	label_clientApMAC, label_clientApName,
	"band", "wifiStandard",
}

var (
	// ── Info / topology ──────────────────────────────────────────────────────

	client_info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_info", Help: "Information about a connected client (value always 1)",
	}, clientBaseLabels)

	// Wired clients — covers both switch-connected and gateway-connected.
	client_wired_info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_wired_info", Help: "Wired client connection topology (value always 1)",
	}, clientWiredLabels)

	// ── Session traffic totals ────────────────────────────────────────────────

	client_traffic_down_bytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_traffic_down_bytes_total", Help: "Total bytes downloaded in current session",
	}, clientBaseLabels)

	client_traffic_up_bytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_traffic_up_bytes_total", Help: "Total bytes uploaded in current session",
	}, clientBaseLabels)

	// ── Real-time activity (Byte/s) ───────────────────────────────────────────
	//
	// Sourced from the API's `activity` (download) and `uploadActivity` (upload) fields.
	// These are instantaneous rates sampled by the controller — treat as gauges,
	// not counters. Wireless clients only report download activity at the top level;
	// uploadActivity is populated for wired clients only.

	client_activity_down_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_activity_down_bytes_per_second",
		Help: "Real-time download rate in bytes per second (from Omada activity field)",
	}, clientBaseLabels)

	client_activity_up_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_activity_up_bytes_per_second",
		Help: "Real-time upload rate in bytes per second (from Omada uploadActivity field; wired clients only)",
	}, clientBaseLabels)

	// ── Uptime ────────────────────────────────────────────────────────────────

	client_uptime_seconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_uptime_seconds", Help: "How long the client has been connected in seconds",
	}, clientBaseLabels)

	// ── Wireless signal / rate ────────────────────────────────────────────────

	client_signal_rssi = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_signal_rssi_dbm", Help: "Wireless signal strength (RSSI) of the primary link in dBm",
	}, clientWirelessLabels)

	client_snr = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_snr_db", Help: "Signal-to-noise ratio of the primary wireless link in dB",
	}, clientWirelessLabels)

	client_rx_rate_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_rx_rate_bps", Help: "Current RX link rate of the primary link in bits per second",
	}, clientWirelessLabels)

	client_tx_rate_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_tx_rate_bps", Help: "Current TX link rate of the primary link in bits per second",
	}, clientWirelessLabels)

	// ── Wired link speed ─────────────────────────────────────────────────────
	//
	// Physical port negotiated speed for switch-connected wired clients.
	// Looked up from PortSpeedMap (built from switch port data at scrape time).
	// Not emitted for gateway-connected clients or switch clients on port 0.

	client_wired_link_speed_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_wired_link_speed_bps",
		Help: "Negotiated physical link speed of the switch port this wired client is connected to, in bits per second",
	}, clientWiredLabels)

	// ── MLO per-link metrics ──────────────────────────────────────────────────

	client_mlo_link_rx_rate_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_mlo_link_rx_rate_bps", Help: "Per-link RX rate for MLO clients in bits per second",
	}, mloLinkLabels)

	client_mlo_link_tx_rate_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_mlo_link_tx_rate_bps", Help: "Per-link TX rate for MLO clients in bits per second",
	}, mloLinkLabels)

	client_mlo_link_snr = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_mlo_link_snr_db", Help: "Per-link SNR for MLO clients in dB",
	}, mloLinkLabels)

	client_mlo_link_rssi = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_mlo_link_rssi_dbm", Help: "Per-link RSSI for MLO clients in dBm",
	}, mloLinkLabels)
)

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
	base[label_clientWifiStd] = c.WifiStandard()
	base[label_clientMLO] = boolToString(c.IsMLO())
	return base
}

func getClientWiredLabels(c Client.Client) prometheus.Labels {
	base := getClientBaseLabels(c)
	base[label_clientConnectDev] = c.ConnectDevice()
	base[label_clientConnectMAC] = c.ConnectDeviceMAC()
	base[label_clientConnectType] = c.ConnectDevType
	base[label_clientSwitchPort] = c.SwitchPortID()
	return base
}

// ExposeClientMetrics resets and repopulates all client metrics.
//
// portSpeeds is a PortSpeedMap built by the scraper from switch port data.
// Pass nil (or an empty map) if port speed data is unavailable — the
// client_wired_link_speed_bps metric simply won't be emitted for that scrape.
func ExposeClientMetrics(clients []Client.Client, portSpeeds Client.PortSpeedMap) {
	client_info.Reset()
	client_wired_info.Reset()
	client_traffic_down_bytes.Reset()
	client_traffic_up_bytes.Reset()
	client_activity_down_bps.Reset()
	client_activity_up_bps.Reset()
	client_uptime_seconds.Reset()
	client_signal_rssi.Reset()
	client_snr.Reset()
	client_rx_rate_bps.Reset()
	client_tx_rate_bps.Reset()
	client_wired_link_speed_bps.Reset()
	client_mlo_link_rx_rate_bps.Reset()
	client_mlo_link_tx_rate_bps.Reset()
	client_mlo_link_snr.Reset()
	client_mlo_link_rssi.Reset()

	for _, c := range clients {
		base := getClientBaseLabels(c)

		client_info.With(base).Set(1)
		client_traffic_down_bytes.With(base).Set(c.TrafficDown)
		client_traffic_up_bytes.With(base).Set(c.TrafficUp)
		client_uptime_seconds.With(base).Set(c.Uptime)

		// Real-time activity — emitted for all client types.
		// Download (activity) is populated by the API for both wireless and wired.
		// Upload (uploadActivity) is populated by the API for wired; wireless clients
		// will report 0 when the controller does not provide an upload rate.
		client_activity_down_bps.With(base).Set(c.Activity)
		client_activity_up_bps.With(base).Set(c.UploadActivity)

		if c.Wireless {
			wl := getClientWirelessLabels(c)
			client_signal_rssi.With(wl).Set(c.SignalRSSI())
			client_snr.With(wl).Set(c.SignalSNR())
			client_rx_rate_bps.With(wl).Set(c.LinkRxRate() * 1000) // kbps → bps
			client_tx_rate_bps.With(wl).Set(c.LinkTxRate() * 1000)

			for _, link := range c.MultiLink {
				mloLabels := prometheus.Labels{
					label_clientMAC:    c.MAC,
					label_clientName:   c.DisplayName(),
					label_clientApMAC:  c.ApMAC,
					label_clientApName: c.ApName,
					"band":             link.RadioBand(),
					"wifiStandard":     link.WifiStandardString(),
				}
				client_mlo_link_rx_rate_bps.With(mloLabels).Set(link.RxRate * 1000)
				client_mlo_link_tx_rate_bps.With(mloLabels).Set(link.TxRate * 1000)
				client_mlo_link_snr.With(mloLabels).Set(float64(link.SNR))
				client_mlo_link_rssi.With(mloLabels).Set(float64(link.RSSI))
			}
		} else {
			// Wired — covers both switch-connected and gateway-connected clients.
			wired := getClientWiredLabels(c)
			client_wired_info.With(wired).Set(1)

			// Port link speed — only available for switch clients with a valid port.
			if key := c.WiredPortSpeedKey(); key != "" && portSpeeds != nil {
				if speed, ok := portSpeeds[key]; ok {
					client_wired_link_speed_bps.With(wired).Set(speed)
				}
			}
		}
	}
}
