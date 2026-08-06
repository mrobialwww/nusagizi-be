package models

// Digunakan bersama oleh ChecklistMilestoneTask dan AssessmentKPSPQuestion
type DevelopmentalDomainEnum string

const (
	DomainGrossMotorSkills  DevelopmentalDomainEnum = "gross_motor_skills"
	DomainFineMotorSkills   DevelopmentalDomainEnum = "fine_motor_skills"
	DomainSpeechAndLanguage DevelopmentalDomainEnum = "speech_and_language"
	DomainSocialization     DevelopmentalDomainEnum = "socialization"
)

type MonthTargetEnum string

const (
	MonthTarget3  MonthTargetEnum = "3"
	MonthTarget6  MonthTargetEnum = "6"
	MonthTarget9  MonthTargetEnum = "9"
	MonthTarget12 MonthTargetEnum = "12"
	MonthTarget15 MonthTargetEnum = "15"
	MonthTarget18 MonthTargetEnum = "18"
	MonthTarget21 MonthTargetEnum = "21"
	MonthTarget24 MonthTargetEnum = "24"
	MonthTarget30 MonthTargetEnum = "30"
	MonthTarget36 MonthTargetEnum = "36"
	MonthTarget42 MonthTargetEnum = "42"
	MonthTarget48 MonthTargetEnum = "48"
	MonthTarget54 MonthTargetEnum = "54"
	MonthTarget60 MonthTargetEnum = "60"
)
