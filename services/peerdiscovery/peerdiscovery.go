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

type PeerDiscovery interface {
	GetPeers() ([]entities.Peer, error)
}

type peerDiscoveryImpl struct {
	config entities.Config
}

func NewPeerDiscovery(config entities.Config) PeerDiscovery {
	return &peerDiscoveryImpl{config: config}
}

func (p *peerDiscoveryImpl) GetPeers() ([]entities.Peer, error) {
	peers := map[entities.PeerID]*peerUnderConstruction{}
	var errorsFound []error
	for key, value := range p.config.List() {
		peerID := derivePeerID(key)
		if peerID != "" {
			_, ok := peers[peerID]
			if !ok {
				peers[peerID] = &peerUnderConstruction{
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

	var validPeers []entities.Peer
	for _, peer := range peers {
		peerEntity, ok := peer.toPeer()
		if ok {
			validPeers = append(validPeers, peerEntity)
		}
	}
	return validPeers, errors.Join(errorsFound...)
}

type peerUnderConstruction struct {
	id       entities.PeerID
	address  string
	lastSeen time.Time
}

func (p peerUnderConstruction) isValid() bool {
	return p.id != "" && p.address != "" && p.lastSeen.Unix() > 0
}

func (p peerUnderConstruction) toPeer() (entities.Peer, bool) {
	if !p.isValid() {
		return entities.Peer{}, false
	}
	return entities.NewPeer(p.id, p.address, p.lastSeen), true
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
