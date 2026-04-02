package Enum

type RadioFrequency int8

const (
	RadioFrequency_2_4_Ghz RadioFrequency = 0
	RadioFrequency_5_0_Ghz RadioFrequency = 1
	RadioFrequency_6_0_Ghz RadioFrequency = 2
)

func (rf RadioFrequency) String() string {
	switch rf {
	case RadioFrequency_2_4_Ghz:
		return "2.4 GHz"
	case RadioFrequency_5_0_Ghz:
		return "5.0 GHz"
	case RadioFrequency_6_0_Ghz:
		return "6.0 GHz"
	default:
		return "invalid"
	}
}
