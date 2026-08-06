package models

type MealTimeEnum string

const (
	MealTimeBreakfast      MealTimeEnum = "breakfast"
	MealTimeLunch          MealTimeEnum = "lunch"
	MealTimeDinner         MealTimeEnum = "dinner"
	MealTimeMorningSnack   MealTimeEnum = "morning_snack"
	MealTimeAfternoonSnack MealTimeEnum = "afternoon_snack"
)
