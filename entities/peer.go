package entities

import "time"

type Peer struct {
	id       PeerID
	address  string
	lastSeen time.Time
}

func NewPeer(id PeerID, address string, lastSeen time.Time) Peer {
	return Peer{
		id:       id,
		address:  address,
		lastSeen: lastSeen,
	}
}

func (p Peer) ID() PeerID {
	return p.id
}

func (p Peer) Address() string {
	return p.address
}

func (p Peer) LastSeen() time.Time {
	return p.lastSeen
}

type PeerID string
