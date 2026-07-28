package models

type MealTimeEnum string

const (
	MealTimeSarapan       MealTimeEnum = "SARAPAN"
	MealTimeMakanSiang    MealTimeEnum = "MAKAN_SIANG"
	MealTimeMakanMalam    MealTimeEnum = "MAKAN_MALAM"
	MealTimeSelinganSiang MealTimeEnum = "SELINGAN_SIANG"
	MealTimeSelinganSore  MealTimeEnum = "SELINGAN_SORE"
)

// IsValid memvalidasi bahwa nilai MealTimeEnum sesuai dengan ENUM yang didefinisikan di database.
func (r MealTimeEnum) IsValid() bool {
	switch r {
	case MealTimeSarapan, MealTimeMakanSiang, MealTimeMakanMalam, MealTimeSelinganSiang, MealTimeSelinganSore:
		return true
	}
	return false
}
