package cmd_test

import (
	"os"
	"testing"

	"github.com/geektheripper/gtf/cmd"
)

func TestPublish(t *testing.T) {
	os.Args = []string{"gtf", "m", "p", "test", "internal/log"}
	cmd.Execute()
}
