package log

import "testing"

func TestDebug(t *testing.T) {
	root.SetHandler(StdoutHandler)
	root.Info("stdout handler")
	root.SetHandler(StderrHandler)
	root.Info("stdinfo handler")
}
