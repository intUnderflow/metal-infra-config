package peerdiscovery

import (
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	"regexp"
	"strconv"
	"time"
)

const (
	_errInvalidPeerLastSeen = "invalid peer last seen value"
)

var (
	_peerIDRegex       = regexp.MustCompile("(.+)\\._internal\\..*")
	_peerAddressRegex  = regexp.MustCompile("(.+)\\._internal\\.peer_address")
	_peerLastSeenRegex = regexp.MustCompile("(.+)\\._internal\\.peer_last_seen")
)

type Peer struct {
	id       entities.PeerID
	address  string
	lastSeen time.Time
}

func (p Peer) ID() entities.PeerID {
	return p.id
}

func (p Peer) Address() string {
	return p.address
}

func (p Peer) LastSeen() time.Time {
	return p.lastSeen
}

func (p Peer) isValid() bool {
	return p.id != "" && p.address != "" && p.lastSeen.Unix() > 0
}

type PeerDiscovery struct {
	config *entities.Config
}

func NewPeerDiscovery(config *entities.Config) *PeerDiscovery {
	return &PeerDiscovery{config: config}
}

func (p *PeerDiscovery) GetPeers() ([]Peer, error) {
	peers := map[entities.PeerID]*Peer{}
	var errorsFound []error
	for key, value := range p.config.List() {
		peerID := derivePeerID(key)
		if peerID != "" {
			_, ok := peers[peerID]
			if !ok {
				peers[peerID] = &Peer{
					id: peerID,
				}
			}
			if isPeerAddressKey(key) {
				peers[peerID].address = value.Value
			} else if isPeerLastSeenKey(key) {
				lastSeen, err := strconv.ParseInt(value.Value, 10, 64)
				if err != nil {
					errorsFound = append(errorsFound, errors.Join(errors.New(_errInvalidPeerLastSeen), err))
					continue
				}
				peers[peerID].lastSeen = time.Unix(lastSeen, 0)
			}
		}
	}

	var validPeers []Peer
	for _, peer := range peers {
		if peer.isValid() {
			validPeers = append(validPeers, *peer)
		}
	}
	return validPeers, errors.Join(errorsFound...)
}

func derivePeerID(key entities.Key) entities.PeerID {
	matches := _peerIDRegex.FindStringSubmatch(string(key))
	if len(matches) != 2 {
		return ""
	}
	return entities.PeerID(matches[1])
}

func isPeerAddressKey(key entities.Key) bool {
	return _peerAddressRegex.MatchString(string(key))
}

func isPeerLastSeenKey(key entities.Key) bool {
	return _peerLastSeenRegex.MatchString(string(key))
}
