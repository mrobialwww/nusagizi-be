package models

import (
	"time"

	"github.com/google/uuid"
)

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

type ChildGrowthAnalysis struct {
	ID                  uuid.UUID        `json:"child_growth_analysis_id" db:"child_growth_analysis_id"`
	ChildGrowthReportID uuid.UUID        `json:"child_growth_report_id" db:"child_growth_report_id"`
	AnalysisType        AnalysisTypeEnum `json:"analysis_type" db:"analysis_type"`
	Percentile          float64          `json:"percentile" db:"percentile"`
	InsightTitle        string           `json:"insight_title" db:"insight_title"`
	InsightDescription  string           `json:"insight_description" db:"insight_description"`
	CreatedAt           time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at" db:"updated_at"`
}
