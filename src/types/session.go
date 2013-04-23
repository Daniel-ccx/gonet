package types

import "time"

const (
	FREE = iota
	ONLINE
	RAID				// being raid
)

type Session struct {
	MQ     chan string
	User   User
	Cities []City

	Status	int
	SESSID	[128]byte		// UNIQUE session ID
	HeartBeat time.Time
}