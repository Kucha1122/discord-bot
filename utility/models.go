package utility

type Boss struct {
	ID             string
	Name           string
	Image          string
	KilledBossesY  string
	KilledPlayersY string
	KilledBosses   string
	KilledPlayers  string
	LastSeen       string
	Introduced     string
	KilledDaysAgo  int32
	IntervalKill   float64
}

type BossB struct {
	Name     string
	Respawns string
	Type     string
}

type CharacterResponse struct {
	Characters struct {
		Data struct {
			Name              string `json:"name"`
			Title             string `json:"title"`
			Sex               string `json:"sex"`
			Vocation          string `json:"vocation"`
			Level             int    `json:"level"`
			AchievementPoints int    `json:"achievement_points"`
			World             string `json:"world"`
			Residence         string `json:"residence"`
			LastLogin         []struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"last_login"`
			AccountStatus string `json:"account_status"`
			Status        string `json:"status"`
		} `json:"data"`
		Achievements []interface{} `json:"achievements"`
		Deaths       []struct {
			Date struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
			Level    int           `json:"level"`
			Reason   string        `json:"reason"`
			Involved []interface{} `json:"involved"`
		} `json:"deaths"`
		AccountInformation struct {
			LoyaltyTitle string `json:"loyalty_title"`
			Created      struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"created"`
		} `json:"account_information"`
		OtherCharacters []struct {
			Name   string `json:"name"`
			World  string `json:"world"`
			Status string `json:"status"`
		} `json:"other_characters"`
	} `json:"characters"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}
