package Enum

// WifiMode maps Omada's wifiMode integer to a human-readable WiFi standard string.
type WifiMode int

const (
	WifiMode_BG       WifiMode = 0 // 802.11b/g
	WifiMode_BGN      WifiMode = 1 // 802.11b/g/n
	WifiMode_A        WifiMode = 2 // 802.11a
	WifiMode_AN       WifiMode = 3 // 802.11a/n
	WifiMode_ANAC     WifiMode = 4 // 802.11a/n/ac (WiFi 5)
	WifiMode_ANACAX   WifiMode = 5 // 802.11a/n/ac/ax (WiFi 6)
	WifiMode_BGNAX    WifiMode = 6 // 802.11b/g/n/ax (WiFi 6)
	WifiMode_BGNAXBE  WifiMode = 7 // 802.11b/g/n/ax/be (WiFi 7 2.4GHz)
	WifiMode_ANACAXBE WifiMode = 8 // 802.11a/n/ac/ax/be (WiFi 7 5GHz)
	WifiMode_AXACAXBE WifiMode = 9 // 802.11ax — band determines 6/6E/7 (I dont think the ai's right i think this is also wifi 7
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
		return "WiFi 6 (802.11ax)"
	case WifiMode_BGNAXBE:
		return "WiFi 7 (802.11be 2.4GHz)"
	case WifiMode_ANACAXBE:
		return "WiFi 7 (802.11be 5GHz)"
	case WifiMode_AXACAXBE:
		return "WiFi 7 (802.11be 6GHz)"
	default:
		return "unknown"
	}
}

// StringWithBand returns the most accurate WiFi standard label given the radioId and
// whether the client has been identified as WiFi 7 (via BE wifiMode on any link).
//
// radioId mapping: 0=2.4GHz, 1=5GHz, 3=6GHz
//
// Key rules:
//   - Any radioId=3 (6GHz): WiFi 6E unless IsWifi7 confirmed → WiFi 7
//   - wifiMode=6 on radioId=3: Omada reports ax on 6GHz as mode 6 → still 6E/7 by band
//   - wifiMode=7/8: confirmed WiFi 7 regardless of band
func (wm WifiMode) StringWithBand(radioID int, isWifi7 bool) string {
	// 6GHz band — label by generation, not wifiMode (Omada reuses mode values across bands)
	if radioID == 3 {
		if isWifi7 {
			return "WiFi 7 (802.11be 6GHz)"
		}
		return "WiFi 6E (802.11ax 6GHz)"
	}

	// Confirmed WiFi 7 on other bands
	if wm == WifiMode_BGNAXBE {
		return "WiFi 7 (802.11be 2.4GHz)"
	}
	if wm == WifiMode_ANACAXBE {
		return "WiFi 7 (802.11be 5GHz)"
	}

	// WiFi 6 ax variants — disambiguate by band
	if wm == WifiMode_BGNAX || wm == WifiMode_ANACAX {
		switch radioID {
		case 0:
			return "WiFi 6 (802.11ax 2.4GHz)"
		case 1:
			return "WiFi 6 (802.11ax 5GHz)"
		}
	}

	return wm.String()
}
