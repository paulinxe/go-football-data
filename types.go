package football_data

type Match struct {
	ID          uint         `json:"id" validate:"required"`
	UTCDate     string       `json:"utcDate" validate:"required"`
	Status      string       `json:"status" validate:"required"`
	HomeTeam    Team         `json:"homeTeam" validate:"required"`
	AwayTeam    Team         `json:"awayTeam" validate:"required"`
	Score       Score        `json:"score" validate:"required"`
	Competition *Competition  `json:"competition" validate:"required"`
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
	Winner string `json:"winner"`
	Duration string `json:"duration"`
	FullTime ScoreTime `json:"fullTime"` // TODO: sadly a 0 makes the required validation fail...
	HalfTime ScoreTime `json:"halfTime"` // TODO: sadly a 0 makes the required validation fail...
}

type ScoreTime struct {
	Home uint `json:"home"` // TODO: sadly a 0 makes the required validation fail...
	Away uint `json:"away"` // TODO: sadly a 0 makes the required validation fail..
}

// TODO: Add referees