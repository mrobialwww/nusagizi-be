package who

import "math"

type LMSRow struct {
	L float64
	M float64
	S float64
}

// LMSTable is keyed by age in months
type LMSTable map[int]LMSRow

// LMSTableFloat is keyed by height/length in cm
type LMSTableFloat map[float64]LMSRow


// CalculateZScore applies the WHO LMS formula.
// x is the raw measurement value (weight in kg, height in cm, BMI, etc.)
func CalculateZScore(x, L, M, S float64) float64 {
	if L != 0 {
		return (math.Pow(x/M, L) - 1) / (L * S)
	}
	return math.Log(x/M) / S
}

// ClassifyZScore returns the title and hardcoded description based on z-score value.
func ClassifyZScore(z float64, childName string) (title, description string) {
	switch {
	case z >= -2 && z <= 2:
		title = "Normal"
		description = "Pertumbuhan " + childName + " berada dalam rentang normal sesuai usianya. Pertahankan pola makan seimbang dan pemantauan rutin untuk mendukung pertumbuhan optimal."
	case (z >= -3 && z < -2) || (z > 2 && z <= 3):
		title = "Berisiko"
		description = "Pertumbuhan " + childName + " berada di rentang berisiko. Konsultasikan dengan dokter atau ahli gizi untuk evaluasi lebih lanjut."
	default:
		title = "Kemungkinan Gangguan"
		description = "Pertumbuhan " + childName + " berada di luar rentang normal dan memerlukan perhatian medis segera. Segera konsultasikan dengan dokter."
	}
	return
}

// ClassifyOverallStatus takes a slice of computed z-scores from multiple WHO indicators
// Severity hierarchy: Gangguan Pertumbuhan > Berisiko > Normal.
func ClassifyOverallStatus(scores []float64) (title, description string) {
	worstLevel := 0 // 0 = Normal, 1 = Berisiko, 2 = Gangguan Pertumbuhan

	for _, z := range scores {
		switch {
		case z < -3 || z > 3:
			if worstLevel < 2 {
				worstLevel = 2
			}
		case (z >= -3 && z < -2) || (z > 2 && z <= 3):
			if worstLevel < 1 {
				worstLevel = 1
			}
		}
	}

	switch worstLevel {
	case 2:
		title = "Gangguan Pertumbuhan"
		description = "Pertumbuhan anak berada di luar rentang normal dan memerlukan perhatian medis segera. Segera konsultasikan dengan dokter."
	case 1:
		title = "Berisiko"
		description = "Pertumbuhan anak berada di rentang berisiko pada salah satu atau lebih aspek pengukuran. Konsultasikan dengan dokter atau ahli gizi untuk evaluasi lebih lanjut."
	default:
		title = "Normal"
		description = "Pertumbuhan anak berada dalam rentang normal di semua aspek. Pertahankan pola makan seimbang dan pemantauan rutin."
	}
	return
}
