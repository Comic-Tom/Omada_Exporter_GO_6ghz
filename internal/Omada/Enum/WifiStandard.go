package Enum

// WifiMode maps Omada's wifiMode integer to a human-readable WiFi standard string.
// Note: Omada reuses some wifiMode values across bands — the band (radioId) provides
// the full context. Mode 9 on 5GHz = WiFi 6, mode 9 on 6GHz = WiFi 6E/7.
type WifiMode int

const (
	WifiMode_BG       WifiMode = 0 // 802.11b/g
	WifiMode_BGN      WifiMode = 1 // 802.11b/g/n
	WifiMode_A        WifiMode = 2 // 802.11a
	WifiMode_AN       WifiMode = 3 // 802.11a/n
	WifiMode_ANAC     WifiMode = 4 // 802.11a/n/ac (WiFi 5)
	WifiMode_ANACAX   WifiMode = 5 // 802.11a/n/ac/ax (WiFi 6)
	WifiMode_BGNAX    WifiMode = 6 // 802.11b/g/n/ax (WiFi 6 2.4GHz)
	WifiMode_BGNAXBE  WifiMode = 7 // 802.11b/g/n/ax/be (WiFi 7 2.4GHz)
	WifiMode_ANACAXBE WifiMode = 8 // 802.11a/n/ac/ax/be (WiFi 7 5GHz)
	WifiMode_AX       WifiMode = 9 // 802.11ax (WiFi 6/6E/7 — band determines which)
)

func (wm WifiMode) String() string {
	switch wm {
	case WifiMode_BG:
		return "802.11b/g"
	case WifiMode_BGN:
		return "802.11b/g/n"
	case WifiMode_A:
		return "802.11a"
	case WifiMode_AN:
		return "802.11a/n"
	case WifiMode_ANAC:
		return "WiFi 5 (802.11ac)"
	case WifiMode_ANACAX:
		return "WiFi 6 (802.11ax)"
	case WifiMode_BGNAX:
		return "WiFi 6 (802.11ax 2.4GHz)"
	case WifiMode_BGNAXBE:
		return "WiFi 7 (802.11be 2.4GHz)"
	case WifiMode_ANACAXBE:
		return "WiFi 7 (802.11be 5GHz)"
	case WifiMode_AX:
		return "WiFi 6/7 (802.11ax)"
	default:
		return "unknown"
	}
}

// StringWithBand returns a more specific WiFi standard label when the radioId is known.
// radioId: 0=2.4GHz, 1=5GHz, 3=6GHz
func (wm WifiMode) StringWithBand(radioID int) string {
	if wm == WifiMode_AX {
		switch radioID {
		case 0:
			return "WiFi 6 (802.11ax 2.4GHz)"
		case 1:
			return "WiFi 6 (802.11ax 5GHz)"
		case 3:
			return "WiFi 6E/7 (802.11ax 6GHz)"
		}
	}
	return wm.String()
}
