package utils

import (
	"fmt"
	"time"
)

// FormatAgeString returns a string in the format "X tahun Y bulan" (or just months/years if 0)
// This utility is used across the application to consistently format children's ages.
func FormatAgeString(birthDate time.Time) string {
	now := time.Now()
	
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	
	if now.Day() < birthDate.Day() {
		months--
	}
	
	if months < 0 {
		years--
		months += 12
	}
	
	if years < 0 {
		years = 0
		months = 0
	}
	
	if years == 0 {
		return fmt.Sprintf("%d bulan", months)
	}
	if months == 0 {
		return fmt.Sprintf("%d tahun", years)
	}
	return fmt.Sprintf("%d tahun %d bulan", years, months)
}

// CalculateAgeInMonths calculates the number of full months between two dates.
func CalculateAgeInMonths(birthDate, targetDate time.Time) int {
	months := (targetDate.Year()-birthDate.Year())*12 + int(targetDate.Month()-birthDate.Month())
	if targetDate.Day() < birthDate.Day() {
		months--
	}
	if months < 0 {
		return 0
	}
	return months
}
