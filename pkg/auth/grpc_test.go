package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"testing"
)

func Test_GetPeerIDFromGRPCContext_GetsPeerID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCtx := peer.NewContext(ctx, &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{
					{
						&x509.Certificate{
							Subject: pkix.Name{
								CommonName: "2f5e93b1-9a03-4564-8d95-dab30b5db09e.metal-infra-config.local",
							},
						},
					},
				},
			},
		},
	})

	peerID, err := GetPeerIDFromGRPCContext(testCtx)
	require.NoError(t, err)
	require.Equal(t, entities.PeerID("2f5e93b1-9a03-4564-8d95-dab30b5db09e"), peerID)
}

func Test_GetPeerIDFromGRPCContext_WhenNotAGRPCPeerContext_ReturnsError(t *testing.T) {
	_, err := GetPeerIDFromGRPCContext(context.Background())
	require.EqualError(t, err, ErrContextNotAGRPCPeer)
}

type fakeCredentials struct{}

func (_ fakeCredentials) AuthType() string {
	return "fake"
}

func Test_GetPeerIDFromGRPCContext_WhenNotTLSAuthenticated_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCtx := peer.NewContext(ctx, &peer.Peer{
		AuthInfo: fakeCredentials{},
	})

	_, err := GetPeerIDFromGRPCContext(testCtx)
	require.EqualError(t, err, _errNoTLSInfoForPeer)
}

func Test_GetPeerIDFromGRPCContext_WhenNoCertificate_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCtx := peer.NewContext(ctx, &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{
					{
						// First chain intentionally empty
					},
					{
						&x509.Certificate{
							Subject: pkix.Name{
								CommonName: "foobar",
							},
						},
					},
				},
			},
		},
	})

	_, err := GetPeerIDFromGRPCContext(testCtx)
	require.EqualError(t, err, _errNoCertificate)

	testCtx = peer.NewContext(ctx, &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{
					// No chains
				},
			},
		},
	})

	_, err = GetPeerIDFromGRPCContext(testCtx)
	require.EqualError(t, err, _errNoCertificate)
}
