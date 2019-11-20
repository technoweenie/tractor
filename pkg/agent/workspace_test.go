package agent

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	ag, teardown := setup(t)
	defer teardown()

	t.Run("start/stop", func(t *testing.T) {
		ws := ag.Workspace("test")
		require.NotNil(t, ws)
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		ch := readWorkspace(t, ws.Start)
		time.Sleep(time.Second)
		assert.Equal(t, int(StatusAvailable), int(ws.Status))

		ws.Stop()
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		out := strings.TrimSpace(string(<-ch))
		assert.True(t, strings.HasPrefix(out, "pid "))
	})

	t.Run("start/connect/stop", func(t *testing.T) {
		ws := ag.Workspace("test")
		require.NotNil(t, ws)
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		startCh := readWorkspace(t, ws.Start)
		time.Sleep(time.Second)
		assert.Equal(t, int(StatusAvailable), int(ws.Status))

		connCh := readWorkspace(t, ws.Connect)
		time.Sleep(time.Second)

		ws.Stop()
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		startOut := strings.TrimSpace(string(<-startCh))
		assert.True(t, strings.HasPrefix(startOut, "pid "))

		connOut := strings.TrimSpace(string(<-connCh))
		assert.True(t, strings.HasPrefix(connOut, "pid "))
		assert.Equal(t, startOut, connOut)
	})

	t.Run("connect/start/stop", func(t *testing.T) {
		ws := ag.Workspace("test")
		require.NotNil(t, ws)
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		connCh := readWorkspace(t, ws.Connect)
		time.Sleep(time.Second)
		assert.Equal(t, int(StatusAvailable), int(ws.Status))

		startCh := readWorkspace(t, ws.Start)
		time.Sleep(time.Second)
		assert.Equal(t, int(StatusAvailable), int(ws.Status))

		ws.Stop()
		assert.Equal(t, int(StatusPartially), int(ws.Status))

		startOut := strings.TrimSpace(string(<-startCh))
		assert.True(t, strings.HasPrefix(startOut, "pid "))

		connOut := strings.TrimSpace(string(<-connCh))
		assert.True(t, strings.HasPrefix(connOut, "pid "))
		assert.NotEqual(t, startOut, connOut)
	})
}

func readWorkspace(t *testing.T, wsFunc func() (io.ReadCloser, error)) chan []byte {
	ch := make(chan []byte)
	go func() {
		r, err := wsFunc()
		if err != nil {
			t.Error(err)
			return
		}

		out := &bytes.Buffer{}
		by := make([]byte, 10)
		for {
			n, err := r.Read(by)
			if err != nil {
				if err != io.EOF {
					t.Error(err)
				}
				break
			}
			out.Write(by[0:n])
		}
		ch <- out.Bytes()
	}()
	return ch
}
