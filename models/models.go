package models

import (
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ID                 uint
	Name               string
	WebName            string
	Position           string
	Nationality        string
	TotalPoints        int
	TeamID             uint
	Price              int
	ChanceOfNextRound  int
	CostChange         int
	News               string
	PointsPerGame      string
	SelectedByPercent  string
	Minutes            uint
	Goals              uint
	Assists            uint
	CleanSheets        uint
	GoalsConceded      uint
	OwnGoals           uint
	PenaltiesSaved     uint
	PenaltiesMissed    uint
	YellowCards        uint
	RedCards           uint
	Saves              uint
	Bonus              uint
	Bps                int
	Starts             uint
	XG                 float32
	XA                 float32
	XGI                float32
	XGC                float32
	XGper90            float32
	SavesPer90         float32
	XAper90            float32
	XGIper90           float32
	XGCper90           float32
	GoalsConcededPer90 float32
	CleanSheetsPer90   float32
}

type Team struct {
	gorm.Model
	ID                  uint
	Name                string
	ShortName           string
	StrengthOverallHome uint
	StrengthOverallAway uint
	StrengthAttackHome  uint
	StrengthAttackAway  uint
	StrengthDefenceHome uint
	StrengthDefenceAway uint
}
