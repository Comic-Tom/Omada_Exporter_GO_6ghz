package Enum

// WifiMode maps Omada's wifiMode integer to a human-readable WiFi standard string.
//
// Values from official Omada OpenAPI docs, verified against real device data:
//   0: 11a   → 802.11a           (5GHz legacy)
//   1: 11b   → 802.11b           (2.4GHz legacy)
//   2: 11g   → 802.11g           (2.4GHz)
//   3: 11na  → 802.11a/n         (5GHz WiFi 4)
//   4: 11ng  → 802.11b/g/n       (2.4GHz WiFi 4)
//   5: 11ac  → 802.11a/n/ac      (WiFi 5, 5GHz)
//   6: 11axa → 802.11ax          (WiFi 6, 5GHz ax)   ← NOT 6E; confirmed via radioId=1 in real data
//   7: 11axg → 802.11ax          (WiFi 6, 2.4GHz ax)
//   8: 11beg → 802.11be          (WiFi 7, 2.4GHz be) ← observed: radioId=0 in real multiLink data
//   9: 11bea → 802.11be          (WiFi 7, 5GHz/6GHz be) ← observed: radioId=1 or radioId=3 in real data
//
// Note: wifiMode 9 covers both 5GHz and 6GHz WiFi 7 links. Use radioId to distinguish:
//   radioId=1 → 5GHz, radioId=3 → 6GHz. The String() method returns a generic label;
//   use RadioBandFromID() alongside wifiMode for full band detail.
//
// The official docs define multiLink wifiMode as 0-7, but values 8 and 9 are
// observed in practice for 802.11be (WiFi 7) links.
type WifiMode int

const (
	WifiMode_11a   WifiMode = 0 // 802.11a (5GHz legacy)
	WifiMode_11b   WifiMode = 1 // 802.11b (2.4GHz legacy)
	WifiMode_11g   WifiMode = 2 // 802.11g (2.4GHz)
	WifiMode_11na  WifiMode = 3 // 802.11a/n (5GHz WiFi 4)
	WifiMode_11ng  WifiMode = 4 // 802.11b/g/n (2.4GHz WiFi 4)
	WifiMode_11ac  WifiMode = 5 // 802.11a/n/ac (WiFi 5, 5GHz)
	WifiMode_11axa WifiMode = 6 // 802.11ax (WiFi 6, 5GHz) — NOT 6E
	WifiMode_11axg WifiMode = 7 // 802.11ax (WiFi 6, 2.4GHz)
	WifiMode_11beg WifiMode = 8 // 802.11be (WiFi 7, 2.4GHz)
	WifiMode_11bea WifiMode = 9 // 802.11be (WiFi 7, 5GHz or 6GHz — see radioId)
)

func (wm WifiMode) String() string {
	switch wm {
	case WifiMode_11a:
		return "802.11a"
	case WifiMode_11b:
		return "802.11b"
	case WifiMode_11g:
		return "802.11g"
	case WifiMode_11na:
		return "802.11a/n"
	case WifiMode_11ng:
		return "802.11b/g/n (WiFi 4)"
	case WifiMode_11ac:
		return "WiFi 5 (802.11ac)"
	case WifiMode_11axa:
		return "WiFi 6 (802.11ax 5GHz)"
	case WifiMode_11axg:
		return "WiFi 6 (802.11ax 2.4GHz)"
	case WifiMode_11beg:
		return "WiFi 7 (802.11be 2.4GHz)"
	case WifiMode_11bea:
		return "WiFi 7 (802.11be)"
	default:
		return "unknown"
	}
}

// StringWithBand returns a band-aware WiFi standard label using radioId to resolve
// ambiguous wifiMode values:
//
//   - wifiMode 6 (11axa) on radioId 3 (6GHz) → WiFi 6E; otherwise WiFi 6 5GHz.
//     Confirmed via SilverDragon: wifiMode=6, radioId=3, UI band="6 GHz (11ax)".
//   - wifiMode 7 (11axg) on radioId 3 (6GHz) → WiFi 6E (defensive; not yet observed).
//   - wifiMode 9 (11bea) on radioId 1/2 → WiFi 7 5GHz; radioId 3 → WiFi 7 6GHz.
//
// All other modes encode their band directly in the wifiMode value and fall through
// to String().
func (wm WifiMode) StringWithBand(radioID int) string {
	switch wm {
	case WifiMode_11axa, WifiMode_11axg:
		if radioID == 3 {
			return "WiFi 6E (802.11ax 6GHz)"
		}
	case WifiMode_11bea:
		switch radioID {
		case 1, 2:
			return "WiFi 7 (802.11be 5GHz)"
		case 3:
			return "WiFi 7 (802.11be 6GHz)"
		}
	}
	return wm.String()
}
