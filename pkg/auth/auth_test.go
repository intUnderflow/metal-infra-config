package auth

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getPeerID_WhenGivenValidCertificate_ReturnsPeerID(t *testing.T) {
	cert := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "b5c5a6ba-18df-4a14-9037-f174011ebd61.metal-infra-config.local",
		},
	}

	peerID, err := getPeerID(cert)
	require.NoError(t, err)
	require.Equal(t, entities.PeerID("b5c5a6ba-18df-4a14-9037-f174011ebd61"), peerID)
}

func Test_getPeerID_WhenGivenInvalidCertificate_ReturnsError(t *testing.T) {
	cert := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "lucy.sh",
		},
	}

	peerID, err := getPeerID(cert)
	require.EqualError(t, err, _errCertificateNotForValidPeer)
	require.Empty(t, peerID)
}
