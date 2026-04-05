package Prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"omada_exporter_go/internal/Omada/Model/Client"
)

const (
	label_clientMAC        string = "mac"
	label_clientName       string = "clientName"
	label_clientIP         string = "ip"
	label_clientVendor     string = "vendor"
	label_clientDeviceType string = "clientDeviceType"
	label_clientNetwork    string = "networkName"
	label_clientWireless   string = "wireless"
	label_clientVLAN       string = "vlan"
	label_clientSSID       string = "ssid"
	label_clientBand       string = "band"
	label_clientApMAC      string = "apMac"
	label_clientApName     string = "apName"
	label_clientWifiStd    string = "wifiStandard"
	label_clientMLO        string = "mlo"
	label_clientSwitchMAC  string = "switchMac"
	label_clientSwitchName string = "switchName"
	label_clientSwitchPort string = "switchPort"
)

var clientBaseLabels = []string{
	label_clientMAC, label_clientName, label_clientIP,
	label_clientVendor, label_clientDeviceType, label_clientNetwork,
	label_clientWireless, label_clientVLAN,
}

var clientWirelessLabels = append([]string{}, append(clientBaseLabels,
	label_clientSSID, label_clientBand, label_clientApMAC,
	label_clientApName, label_clientWifiStd, label_clientMLO,
)...)

var clientWiredLabels = append([]string{}, append(clientBaseLabels,
	label_clientSwitchMAC, label_clientSwitchName, label_clientSwitchPort,
)...)

var mloLinkLabels = []string{
	label_clientMAC, label_clientName, label_clientApMAC,
	label_clientApName, "band", "wifiStandard",
}

var (
	client_info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_info", Help: "Information about a connected client (value always 1)",
	}, clientBaseLabels)

	client_traffic_down_bytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_traffic_down_bytes_total", Help: "Total bytes downloaded by the client in the current session",
	}, clientBaseLabels)

	client_traffic_up_bytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_traffic_up_bytes_total", Help: "Total bytes uploaded by the client in the current session",
	}, clientBaseLabels)

	client_uptime_seconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_uptime_seconds", Help: "How long the client has been connected in seconds",
	}, clientBaseLabels)

	// Wireless — uses primary link data (best signal link for MLO)
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

	// MLO per-link breakdown
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

	// Wired
	client_wired_link_speed_bps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_wired_link_speed_bps", Help: "Wired link speed placeholder — topology visible in labels",
	}, clientWiredLabels)
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
	base[label_clientSwitchMAC] = c.SwitchMAC
	base[label_clientSwitchName] = c.SwitchName
	base[label_clientSwitchPort] = c.SwitchPortID()
	return base
}

func ExposeClientMetrics(clients []Client.Client) {
	client_info.Reset()
	client_traffic_down_bytes.Reset()
	client_traffic_up_bytes.Reset()
	client_uptime_seconds.Reset()
	client_signal_rssi.Reset()
	client_snr.Reset()
	client_rx_rate_bps.Reset()
	client_tx_rate_bps.Reset()
	client_mlo_link_rx_rate_bps.Reset()
	client_mlo_link_tx_rate_bps.Reset()
	client_mlo_link_snr.Reset()
	client_mlo_link_rssi.Reset()
	client_wired_link_speed_bps.Reset()

	for _, c := range clients {
		base := getClientBaseLabels(c)
		client_info.With(base).Set(1)
		client_traffic_down_bytes.With(base).Set(c.TrafficDown)
		client_traffic_up_bytes.With(base).Set(c.TrafficUp)
		client_uptime_seconds.With(base).Set(c.Uptime)

		if c.Wireless {
			wl := getClientWirelessLabels(c)
			// Use primary link RSSI/SNR/rates (best signal link for MLO)
			client_signal_rssi.With(wl).Set(c.SignalRSSI())
			client_snr.With(wl).Set(c.SignalSNR())
			client_rx_rate_bps.With(wl).Set(c.LinkRxRate() * 1000) // kbps → bps
			client_tx_rate_bps.With(wl).Set(c.LinkTxRate() * 1000)

			// Per-link MLO breakdown — emit all links always
			for _, link := range c.MultiLink {
				mloLabels := prometheus.Labels{
					label_clientMAC:    c.MAC,
					label_clientName:   c.DisplayName(),
					label_clientApMAC:  c.ApMAC,
					label_clientApName: c.ApName,
					"band":             link.RadioBand(),
					// Use band-aware standard string for accurate labelling
					// We pass c.IsMLO() or a similar check to satisfy the second boolean argument
"wifiStandard": link.WifiMode.StringWithBand(link.RadioID, c.IsMLO()),
				}
				client_mlo_link_rx_rate_bps.With(mloLabels).Set(link.RxRate * 1000)
				client_mlo_link_tx_rate_bps.With(mloLabels).Set(link.TxRate * 1000)
				client_mlo_link_snr.With(mloLabels).Set(float64(link.SNR))
				client_mlo_link_rssi.With(mloLabels).Set(float64(link.RSSI))
			}
		} else {
			client_wired_link_speed_bps.With(getClientWiredLabels(c)).Set(0)
		}
	}
}
