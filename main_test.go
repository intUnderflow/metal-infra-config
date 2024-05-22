package main

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/goleak"
	"testing"
)

func TestFX(t *testing.T) {
	require.NoError(t, fx.ValidateApp(opts()))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
