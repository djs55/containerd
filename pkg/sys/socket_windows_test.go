package sys

import (
	"net"
	"path/filepath"
	"testing"
)

func TestIsNamedPipePath(t *testing.T) {
	testcases := []struct {
		name string
		path string
		want bool
	}{
		{name: "forward slash pipe", path: "//./pipe/containerd", want: true},
		{name: "backslash pipe", path: `\\.\pipe\containerd`, want: true},
		{name: "unix socket path", path: "/tmp/test.sock", want: false},
		{name: "windows fs path", path: `C:\Users\test\docker.sock`, want: false},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isNamedPipePath(tc.path); got != tc.want {
				t.Errorf("isNamedPipePath(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestGetLocalListenerUnixSocket(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	l, err := GetLocalListener(sockPath, 0, 0)
	if err != nil {
		t.Fatalf("GetLocalListener(%q) failed: %v", sockPath, err)
	}
	defer l.Close()

	// Verify we can connect to it
	done := make(chan error, 1)
	go func() {
		conn, err := l.Accept()
		if err == nil {
			conn.Close()
		}
		done <- err
	}()

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	conn.Close()

	if err := <-done; err != nil {
		t.Fatalf("Accept failed: %v", err)
	}
}
