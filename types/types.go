package types

type MatchesList struct {
	Matches []Match `json:"matches" validate:"required,dive"`
}

type Match struct {
	ID          uint         `json:"id" validate:"required"`
	UTCDate     string       `json:"utcDate" validate:"required"`
	Status      string       `json:"status" validate:"required"`
	HomeTeam    Team         `json:"homeTeam" validate:"required"`
	AwayTeam    Team         `json:"awayTeam" validate:"required"`
	Score       Score        `json:"score" validate:"required"`
	Competition Competition  `json:"competition" validate:"required"`
}

type Competition struct {
	ID uint `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
	Type string `json:"type" validate:"required"`
	Emblem string `json:"emblem" validate:"required"`
}

type Team struct {
	ID uint `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	ShortName string `json:"shortName" validate:"required"`
	TLA string `json:"tla" validate:"required"`
	Crest string `json:"crest" validate:"required"`
}

type Score struct {
	Winner *string `json:"winner"`
	Duration string `json:"duration" validate:"required"`
	FullTime ScoreTime `json:"fullTime"`
	HalfTime ScoreTime `json:"halfTime"`
}

type ScoreTime struct {
	Home *uint `json:"home"`
	Away *uint `json:"away"`
}

// TODO: Add referees