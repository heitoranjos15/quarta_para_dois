package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger_ProductionIsJSON(t *testing.T) {
	var buf bytes.Buffer
	newLogger("production", &buf).Info("hi")

	var m map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &m))
	assert.Equal(t, "hi", m["msg"])
}

func TestNewLogger_DevelopmentIsText(t *testing.T) {
	var buf bytes.Buffer
	newLogger("development", &buf).Info("hi")

	var m map[string]any
	assert.Error(t, json.Unmarshal(buf.Bytes(), &m))
	assert.Contains(t, buf.String(), "msg=hi")
}

func TestNewServer_Timeouts(t *testing.T) {
	srv := newServer(":8080", http.NewServeMux())

	assert.Equal(t, ":8080", srv.Addr)
	assert.Equal(t, 10*time.Second, srv.ReadHeaderTimeout)
	assert.Equal(t, 15*time.Second, srv.ReadTimeout)
	assert.Equal(t, 60*time.Second, srv.WriteTimeout)
	assert.Equal(t, 120*time.Second, srv.IdleTimeout)
}
