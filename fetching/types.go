package fetching

type FetchedPlayer struct {
	ID                 int     `json:"id"`
	FirstName          string  `json:"first_name"`
	SecondName         string  `json:"second_name"`
	WebName            string  `json:"web_name"`
	ElementType        int     `json:"element_type"`
	TotalPoints        int     `json:"total_points"`
	Team               int     `json:"team"`
	NowCost            int     `json:"now_cost"`
	ChanceOfNextRound  int     `json:"chance_of_playing_next_round"`
	CostChange         int     `json:"cost_change_start"`
	News               string  `json:"news"`
	Minutes            int     `json:"minutes"`
	GoalsScored        int     `json:"goals_scored"`
	Assists            int     `json:"assists"`
	CleanSheets        int     `json:"clean_sheets"`
	GoalsConceded      int     `json:"goals_conceded"`
	Saves              int     `json:"saves"`
	PenaltiesSaved     int     `json:"penalties_saved"`
	PenaltiesMissed    int     `json:"penalties_missed"`
	YellowCards        int     `json:"yellow_cards"`
	RedCards           int     `json:"red_cards"`
	Form               string  `json:"form"`
	PointsPerGame      string  `json:"points_per_game"`
	SelectedByPercent  string  `json:"selected_by_percent"`
	TransfersIn        int     `json:"transfers_in"`
	TransfersOut       int     `json:"transfers_out"`
	Bonus              int     `json:"bonus"`
	Bps                int     `json:"bps"`
	Starts             int     `json:"starts"`
	XG                 string  `json:"expected_goals"`
	XA                 string  `json:"expected_assists"`
	XGI                string  `json:"expected_goal_involvements"`
	XGC                string  `json:"expected_goals_conceded"`
	XGper90            float32 `json:"expected_goals_per_90"`
	SavesPer90         float32 `json:"saves_per_90"`
	XAper90            float32 `json:"expected_assists_per_90"`
	XGIper90           float32 `json:"expected_goal_involvements_per_90"`
	XGCper90           float32 `json:"expected_goals_conceded_per_90"`
	GoalsConcededPer90 float32 `json:"goals_conceded_per_90"`
	CleanSheetsPer90   float32 `json:"clean_sheets_per_90"`
}

type FetchedTeam struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	ShortName           string `json:"short_name"`
	StrengthOverallHome int    `json:"strength_overall_home"`
	StrengthOverallAway int    `json:"strength_overall_away"`
	StrengthAttackHome  int    `json:"strength_attack_home"`
	StrengthAttackAway  int    `json:"strength_attack_away"`
	StrengthDefenceHome int    `json:"strength_defence_home"`
	StrengthDefenceAway int    `json:"strength_defence_away"`
}

type TeamsAndPlayersData struct {
	Players []FetchedPlayer `json:"elements"`
	Teams   []FetchedTeam   `json:"teams"`
}
