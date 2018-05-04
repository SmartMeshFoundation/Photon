package log

import "testing"

func TestDebug(t *testing.T) {
	root.SetHandler(StdoutHandler)
	root.Info("stdout handler")
	root.SetHandler(StderrHandler)
	root.Info("stdinfo handler")
	l := New("id", 30, "conn", "ok")
	l.Info("from id")
	l2 := l.New("submodule", "test")
	l2.Info("from l2")
}
