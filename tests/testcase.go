package tests

import (
	"bytes"
	"log"
	"os"
)

type TestCase struct {
	input         string
	shouldSucceed bool
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
