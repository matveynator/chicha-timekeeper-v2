package Data

type RawData  struct {
	TagId                       string
	DiscoveryUnixTime           int64 
	ReaderIP                    string
	Antenna                     uint8
	ProxyIP											string
}

type RawXMLData  struct {
	TagId                       string
	DiscoveryTime           		string
	Antenna                     uint8
	ReaderIP										string
}


type Lap struct {
	Id                          int64
	SportsmanId                 string
	TagId                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	AverageResultsCount					int
	RaceId                      uint
	RacePosition                uint
	TimeBehindTheLeader         int64
	LapNumber                   int
	LapTime                     int64
	LapPosition                 uint
	LapIsCurrent                bool
	LapIsStrange                bool
	RaceFinished                bool
	BestLapTime                 int64
	BestLapNumber               int
	RaceTotalTime               int64
	BetterOrWorseLapTime        int64
}


