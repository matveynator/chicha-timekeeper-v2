package Data

type RawData  struct {
	TagID                       string
	DiscoveryUnixTime           int64
	Antenna                     uint8
	ReaderIP                    string
	ProxyIP											string
}

type Lap struct {
	ID                          int64
	SportsmanID                 string
	TagID                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	AverageResultsCount					int
	Antenna											int
	AntennaIP										string

	UpdatedAt                   int64
	RaceID                      uint
	PracticeID                  uint
	RacePosition                uint
	TimeBehindTheLeader         int64
	LapNumber                   int
	LapTime                     int64
	LapPosition                 uint
	LapIsCurrent                bool
	LapIsStrange                bool
	StageFinished								bool
	RaceFinished                bool
	BestLapTime                 int64
	BestLapNumber               int
	BestLapPosition             uint
	RaceTotalTime               int64
	BetterOrWorseLapTime        int64
}


