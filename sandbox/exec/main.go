package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/google/uuid"
)

func main() {
	sessionId := uuid.New().String()
	slog.Info("session id", slog.String("sessionId", sessionId))
	cmd := exec.Command(
		"claude",
		"-p", "system prompt",
		"--session-id", sessionId,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("error starting cmd", slog.Any("err", err))
		panic(err)
	}

	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
}
