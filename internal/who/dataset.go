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

// Measurement is a minimal input type for who analysis functions.
type Measurement struct {
	WeightKg            float64
	HeightCm            float64
	HeadCircumferenceCm float64
	AgeMonths           int
}

// AnalysisPoint is the output of a single plotted measurement.
type AnalysisPoint struct {
	X float64
	Y float64
}

// AnalysisMessage carries the Z-Score interpretation for the latest measurement.
type AnalysisMessage struct {
	Title       string
	Description string
}

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
		title = "Sangat Buruk"
		description = "Pertumbuhan " + childName + " berada di luar rentang normal dan memerlukan perhatian medis segera. Segera konsultasikan dengan dokter."
	}
	return
}

// ComputeGrowthStatus calculates Z-Scores for all 5 WHO indicators and returns the overall status.
func ComputeGrowthStatus(gender string, ageMonths int, weightKg, heightCm, headCm float64) (status, description string) {
	var scores []float64
	var lms LMSRow
	var found bool

	isMale := gender == "male"

	if weightKg > 0 {
		if isMale {
			lms, found = WFABoys[ageMonths]
		} else {
			lms, found = WFAGirls[ageMonths]
		}
		if found {
			scores = append(scores, CalculateZScore(weightKg, lms.L, lms.M, lms.S))
		}
	}

	if heightCm > 0 {
		if isMale {
			lms, found = LHFABoys[ageMonths]
		} else {
			lms, found = LHFAGirls[ageMonths]
		}
		if found {
			scores = append(scores, CalculateZScore(heightCm, lms.L, lms.M, lms.S))
		}
	}

	if weightKg > 0 && heightCm > 0 {
		roundedH := math.Round(heightCm*2) / 2 // round to nearest 0.5 cm
		if ageMonths < 24 {
			if isMale {
				lms, found = WFLBoys[roundedH]
			} else {
				lms, found = WFLGirls[roundedH]
			}
		} else {
			if isMale {
				lms, found = WFHBoys[roundedH]
			} else {
				lms, found = WFHGirls[roundedH]
			}
		}
		if found {
			scores = append(scores, CalculateZScore(weightKg, lms.L, lms.M, lms.S))
		}
	}

	if weightKg > 0 && heightCm > 0 {
		heightM := heightCm / 100
		bmi := weightKg / (heightM * heightM)
		if isMale {
			lms, found = BFABoys[ageMonths]
		} else {
			lms, found = BFAGirls[ageMonths]
		}
		if found {
			scores = append(scores, CalculateZScore(bmi, lms.L, lms.M, lms.S))
		}
	}

	if headCm > 0 {
		if isMale {
			lms, found = HCFABoys[ageMonths]
		} else {
			lms, found = HCFAGirls[ageMonths]
		}
		if found {
			scores = append(scores, CalculateZScore(headCm, lms.L, lms.M, lms.S))
		}
	}

	if len(scores) == 0 {
		return "Tidak Tersedia", "Data pengukuran tidak mencukupi untuk menghitung status pertumbuhan."
	}

	// 2. Classify based on the worst Z-Score
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
		status = "Sangat Buruk"
		description = "Pertumbuhan anak berada di luar rentang normal dan memerlukan perhatian medis segera. Segera konsultasikan dengan dokter."
	case 1:
		status = "Berisiko"
		description = "Pertumbuhan anak berada di rentang berisiko pada salah satu atau lebih aspek pengukuran. Konsultasikan dengan dokter atau ahli gizi untuk evaluasi lebih lanjut."
	default:
		status = "Normal"
		description = "Pertumbuhan anak berada dalam rentang normal di semua aspek. Pertahankan pola makan seimbang dan pemantauan rutin."
	}
	return
}

// BuildAnalysisPoints transforms raw measurements into chart coordinates (X,Y) and extracts WHO Z-Score status from the latest data.
func BuildAnalysisPoints(gender, analysisType string, measurements []Measurement) ([]AnalysisPoint, *AnalysisMessage) {
	points := make([]AnalysisPoint, 0, len(measurements))
	var message *AnalysisMessage

	for i, m := range measurements {
		isLatest := (i == len(measurements)-1)
		x, y, lms, found := resolveIndicator(gender, analysisType, m)

		points = append(points, AnalysisPoint{X: x, Y: y})

		if isLatest && found {
			z := CalculateZScore(y, lms.L, lms.M, lms.S)
			title, desc := ClassifyZScore(z, "anak")
			message = &AnalysisMessage{Title: title, Description: desc}
		}
	}

	return points, message
}

// resolveIndicator maps a measurement to its specific WHO reference table (LMS) and axis values based on gender and indicator type.
func resolveIndicator(gender, analysisType string, m Measurement) (x, y float64, lms LMSRow, found bool) {
	isMale := gender == "male"

	switch analysisType {
	case "weight_for_age":
		x = float64(m.AgeMonths)
		y = m.WeightKg
		if isMale {
			lms, found = WFABoys[m.AgeMonths]
		} else {
			lms, found = WFAGirls[m.AgeMonths]
		}

	case "height_for_age":
		x = float64(m.AgeMonths)
		y = m.HeightCm
		if isMale {
			lms, found = LHFABoys[m.AgeMonths]
		} else {
			lms, found = LHFAGirls[m.AgeMonths]
		}

	case "weight_for_height":
		roundedH := math.Round(m.HeightCm*2) / 2
		x = roundedH
		y = m.WeightKg
		if m.AgeMonths < 24 {
			if isMale {
				lms, found = WFLBoys[roundedH]
			} else {
				lms, found = WFLGirls[roundedH]
			}
		} else {
			if isMale {
				lms, found = WFHBoys[roundedH]
			} else {
				lms, found = WFHGirls[roundedH]
			}
		}

	case "bmi_for_age":
		x = float64(m.AgeMonths)
		if m.HeightCm > 0 {
			heightM := m.HeightCm / 100.0
			y = m.WeightKg / (heightM * heightM)
			if isMale {
				lms, found = BFABoys[m.AgeMonths]
			} else {
				lms, found = BFAGirls[m.AgeMonths]
			}
		}

	case "head_circumference_for_age":
		x = float64(m.AgeMonths)
		y = m.HeadCircumferenceCm
		if isMale {
			lms, found = HCFABoys[m.AgeMonths]
		} else {
			lms, found = HCFAGirls[m.AgeMonths]
		}
	}

	return
}
