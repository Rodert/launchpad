package common

type (
	Weights struct {
		BehavioralRisk       float64 `json:"behavioralRisk"`
		EntityRisk           float64 `json:"entityRisk"`
		AnomalousAssociation float64 `json:"anomalousAssociation"`
	}
	AlertLevels struct {
		HighThreshold   int `json:"highThreshold"`
		MediumThreshold int `json:"mediumThreshold"`
	}
)

type (
	Policy struct {
		Weights     Weights     `json:"weights"`
		AlertLevels AlertLevels `json:"alertLevels"`
	}
)
