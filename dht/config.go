package dht

type config struct {
	Address string
	Ip      string
	Port    int

	ReadBuffer int

	PrimeNodes []string
}

func DefualtConfig() *config {
	return &config{
		Address: "localhost:6881",
		Ip:      "localhost",
		Port:    6881,

		ReadBuffer: 1024,

		PrimeNodes: []string{
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
	}
}
