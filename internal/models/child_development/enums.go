package models

// Digunakan bersama oleh ChecklistMilestoneTask dan AssessmentKPSPQuestion
type NerveNameEnum string

const (
	NerveGrossMotorSkills  NerveNameEnum = "GROSS_MOTOR_SKILLS"
	NerveFineMotorSkills   NerveNameEnum = "FINE_MOTOR_SKILLS"
	NerveSpeechAndLanguage NerveNameEnum = "SPEECH_AND_LANGUAGE"
	NerveSocialization     NerveNameEnum = "SOCIALIZATION"
)

// IsValid memastikan nilai NerveNameEnum valid sesuai ENUM database
func (n NerveNameEnum) IsValid() bool {
	switch n {
	case NerveGrossMotorSkills, NerveFineMotorSkills, NerveSpeechAndLanguage, NerveSocialization:
		return true
	}
	return false
}

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

// IsValid memastikan nilai MonthTargetEnum yang di-input valid sesuai ENUM database
func (m MonthTargetEnum) IsValid() bool {
	switch m {
	case MonthTarget3, MonthTarget6, MonthTarget9, MonthTarget12, MonthTarget15,
		MonthTarget18, MonthTarget21, MonthTarget24, MonthTarget30, MonthTarget36,
		MonthTarget42, MonthTarget48, MonthTarget54, MonthTarget60:
		return true
	}
	return false
}
