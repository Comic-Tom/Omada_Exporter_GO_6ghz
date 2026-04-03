package Enum

// WifiMode maps Omada's wifiMode integer to a human-readable WiFi standard string.
// Values derived from empirical observation of Omada API behavior with WiFi 6/7 devices.
type WifiMode int

const (
	WifiMode_A   WifiMode = 0 // 802.11a
	WifiMode_B   WifiMode = 1 // 802.11b
	WifiMode_G   WifiMode = 2 // 802.11g
	WifiMode_NA  WifiMode = 3 // 802.11n (5GHz)
	WifiMode_NG  WifiMode = 4 // 802.11n (2.4GHz)
	WifiMode_AC  WifiMode = 5 // 802.11ac (WiFi 5)
	WifiMode_AXA WifiMode = 6 // 802.11ax (WiFi 6 5/6GHz)
	WifiMode_AXG WifiMode = 7 // 802.11ax (WiFi 6 2.4GHz)
	WifiMode_BEG WifiMode = 8 // 802.11be (WiFi 7 2.4GHz)
	WifiMode_BEA WifiMode = 9 // 802.11be (WiFi 7 5/6GHz)
)

func (wm WifiMode) String() string {
	switch wm {
	case WifiMode_A:
		return "802.11a"
	case WifiMode_B:
		return "802.11b"
	case WifiMode_G:
		return "802.11g"
	case WifiMode_NA:
		return "802.11n (5GHz)"
	case WifiMode_NG:
		return "802.11n (2.4GHz)"
	case WifiMode_AC:
		return "WiFi 5 (802.11ac)"
	case WifiMode_AXA:
		return "WiFi 6 (802.11ax 5/6GHz)"
	case WifiMode_AXG:
		return "WiFi 6 (802.11ax 2.4GHz)"
	case WifiMode_BEG:
		return "WiFi 7 (802.11be 2.4GHz)"
	case WifiMode_BEA:
		return "WiFi 7 (802.11be 5/6GHz)"
	default:
		return "unknown"
	}
}
