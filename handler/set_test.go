package handler

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/pkg/auth"
	"github.com/intunderflow/metal-infra-config/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"testing"
)

func Test_Set_WhenNoKey_ReturnsError(t *testing.T) {
	_, err := NewHandler(entities.NewConfig()).Set(context.Background(), &proto.SetRequest{
		Key:     "",
		Value:   "foobar",
		Version: 2,
	})

	require.EqualError(t, err, _errKeyRequired)
}

func Test_Set_WhenNoValue_ReturnsError(t *testing.T) {
	_, err := NewHandler(entities.NewConfig()).Set(context.Background(), &proto.SetRequest{
		Key:     "foobar",
		Value:   "",
		Version: 2,
	})

	require.EqualError(t, err, _errValueRequired)
}

func Test_Set_WhenNoVersion_ReturnsError(t *testing.T) {
	_, err := NewHandler(entities.NewConfig()).Set(context.Background(), &proto.SetRequest{
		Key:     "foobar",
		Value:   "foobar",
		Version: 0,
	})

	require.EqualError(t, err, _errVersionRequired)
}

func Test_Set_WhenBadPeer_ReturnsError(t *testing.T) {
	_, err := NewHandler(entities.NewConfig()).Set(context.Background(), &proto.SetRequest{
		Key:     "foo",
		Value:   "bar",
		Version: 2,
	})

	require.EqualError(t, err, auth.ErrContextNotAGRPCPeer)
}

func Test_Set_SetsValueInConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := entities.NewConfig()
	handler := NewHandler(config)

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

	_, err := handler.Set(testCtx, &proto.SetRequest{
		Key:     "lu",
		Value:   "cy",
		Version: 1,
	})
	require.NoError(t, err)

	value, err := config.GetWithVersion("2f5e93b1-9a03-4564-8d95-dab30b5db09e.lu")
	require.NoError(t, err)
	require.Equal(t, "cy", value.Value)
}

func Test_Set_WhenVersionAlreadyHigher_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := entities.NewConfig()
	require.NoError(t, config.SetWithVersion("2f5e93b1-9a03-4564-8d95-dab30b5db09e.foo", "bar", 10))
	handler := NewHandler(config)

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

	_, err := handler.Set(testCtx, &proto.SetRequest{
		Key:     "foo",
		Value:   "baz",
		Version: 1,
	})
	require.EqualError(t, err, entities.ErrNewValueIsOlderThanCurrentValue)
}
