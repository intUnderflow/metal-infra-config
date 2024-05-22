package auth

import (
	"crypto/x509"
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	"regexp"
)

const (
	_errCertificateNotForValidPeer = "certificate is not for a valid metal-infra-config peer"
)

var (
	_peerIDSubjectRegex = regexp.MustCompile("(.+)\\.metal-infra-config\\.local")
)

func getPeerID(certificate *x509.Certificate) (entities.PeerID, error) {
	matches := _peerIDSubjectRegex.FindStringSubmatch(certificate.Subject.CommonName)
	if len(matches) != 2 {
		return "", errors.New(_errCertificateNotForValidPeer)
	}

	return entities.PeerID(matches[1]), nil
}
