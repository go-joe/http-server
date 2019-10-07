package joehttp

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/go-joe/joe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func joeConf(t *testing.T) *joe.Config {
	joeConf := new(joe.Config)
	require.NoError(t, joe.WithLogger(zaptest.NewLogger(t)).Apply(joeConf))
	return joeConf
}

func TestWithLogger(t *testing.T) {
	logger := zaptest.NewLogger(t)
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithLogger(logger),
	})

	require.NoError(t, err)
	assert.Equal(t, logger, conf.logger)
}

func TestWithTLS(t *testing.T) {
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithTLS("foo.cert", "foo.key"),
	})

	require.NoError(t, err)
	assert.Equal(t, "foo.cert", conf.certFile)
	assert.Equal(t, "foo.key", conf.keyFile)
}

func TestWithTLSConfig(t *testing.T) {
	tlsConf := &tls.Config{ServerName: "foo"}
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithTLSConfig(tlsConf),
	})

	require.NoError(t, err)
	assert.Equal(t, tlsConf, conf.tlsConf)
}

func TestWithTimeouts(t *testing.T) {
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithTimeouts(42 * time.Second),
	})

	require.NoError(t, err)
	assert.EqualValues(t, 42*time.Second, conf.readTimeout)
	assert.EqualValues(t, 42*time.Second, conf.writeTimeout)
}

func TestWithReadTimeout(t *testing.T) {
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithReadTimeout(42 * time.Second),
	})

	require.NoError(t, err)
	assert.EqualValues(t, 42*time.Second, conf.readTimeout)
	assert.EqualValues(t, 0, conf.writeTimeout)
}

func TestWithWriteTimeout(t *testing.T) {
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithWriteTimeout(42 * time.Second),
	})

	require.NoError(t, err)
	assert.EqualValues(t, 0, conf.readTimeout)
	assert.EqualValues(t, 42*time.Second, conf.writeTimeout)
}

func TestWithTrustedHeader(t *testing.T) {
	conf, err := newConf("localhost:0", joeConf(t), []Option{
		WithTrustedHeader("x-real-ip"),
	})
	require.NoError(t, err)
	assert.EqualValues(t, "x-real-ip", conf.trustedHeader)
}
