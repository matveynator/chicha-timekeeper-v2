package Data

type RawData  struct {
	TagId                       string
	DiscoveryUnixTime           int64 
	ReaderIP                    string
	Antenna                     uint
	ProxyIP											string
}

type RawXMLData  struct {
	TagId                       string
	DiscoveryTime           		string
	Antenna                     uint
	ReaderIP										string
}


type Lap struct {
	Id                          int64
	SportsmanId                 string
	TagId                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	AverageResultsCount					uint
	RaceId                      uint
	RacePosition                uint
	TimeBehindTheLeader         int64
	LapNumber                   uint
	LapTime                     int64
	LapPosition                 uint
	LapIsCurrent                bool
	LapIsStrange                bool
	RaceFinished                bool
	BestLapTime                 int64
	BestLapNumber               uint
	RaceTotalTime               int64
	BetterOrWorseLapTime        int64
}


