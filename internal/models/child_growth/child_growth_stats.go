package models



type AnalysisTypeEnum string

const (
	AnalysisTypeWeightForAge            AnalysisTypeEnum = "WEIGHT_FOR_AGE"
	AnalysisTypeHeightForAge            AnalysisTypeEnum = "HEIGHT_FOR_AGE"
	AnalysisTypeWeightForHeight         AnalysisTypeEnum = "WEIGHT_FOR_HEIGHT"
	AnalysisTypeBMIForAge               AnalysisTypeEnum = "BMI_FOR_AGE"
	AnalysisTypeHeadCircumferenceForAge AnalysisTypeEnum = "HEAD_CIRCUMFERENCE_FOR_AGE"
)

// IsValid memvalidasi bahwa nilai AnalysisTypeEnum sesuai dengan ENUM yang didefinisikan di database.
func (a AnalysisTypeEnum) IsValid() bool {
	switch a {
	case AnalysisTypeWeightForAge, AnalysisTypeHeightForAge, AnalysisTypeWeightForHeight,
		AnalysisTypeBMIForAge, AnalysisTypeHeadCircumferenceForAge:
		return true
	}
	return false
}
