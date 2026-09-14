package models

type AnalysisTypeEnum string

const (
	AnalysisTypeWeightForAge            AnalysisTypeEnum = "WEIGHT_FOR_AGE"
	AnalysisTypeHeightForAge            AnalysisTypeEnum = "HEIGHT_FOR_AGE"
	AnalysisTypeWeightForHeight         AnalysisTypeEnum = "WEIGHT_FOR_HEIGHT"
	AnalysisTypeBMIForAge               AnalysisTypeEnum = "BMI_FOR_AGE"
	AnalysisTypeHeadCircumferenceForAge AnalysisTypeEnum = "HEAD_CIRCUMFERENCE_FOR_AGE"
)
