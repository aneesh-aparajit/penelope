package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type process struct {
	sessionId  string
	workflowId string
	message    string
	logger     *slog.Logger
	cmd        *exec.Cmd
	stdout     bytes.Buffer
	stderr     bytes.Buffer
}

func main() {
	root := &cobra.Command{Use: "penelope"}

	fi, err := os.OpenFile(".logs.jsonl", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error opening file", slog.Any("err", err))
		panic(err)
	}
	proc := &process{
		logger: slog.New(slog.NewJSONHandler(fi, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}

	start := &cobra.Command{
		Use:   "start",
		Short: "start a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			if proc.sessionId == "" {
				proc.sessionId = uuid.NewString()
			}
			if proc.workflowId == "" {
				proc.workflowId = uuid.NewString()
			}
			proc.logger.Info("penelope start",
				slog.String("sessionId", proc.sessionId),
				slog.String("workflowId", proc.workflowId),
			)
			proc.cmd = exec.Command(
				"claude", "-p", "system prompt",
				"--session-id", proc.sessionId,
			)
			proc.cmd.Env = append(proc.cmd.Environ(),
				fmt.Sprintf("PENELOPE_WF_ID=%v", proc.workflowId),
				fmt.Sprintf("PENELOPE_SESSION_ID=%v", proc.sessionId),
			)
			proc.cmd.Stderr = &proc.stderr
			proc.cmd.Stdout = &proc.stdout

			if err := proc.cmd.Start(); err != nil {
				proc.logger.Error("error starting claude session",
					slog.Any("err", err))
				return err
			}
			return nil
		},
	}

	send := &cobra.Command{
		Use:   "send",
		Short: "send a message to a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			proc.logger.Info("penelope send",
				slog.String("sessionId", proc.sessionId),
				slog.String("workflowId", proc.workflowId),
			)
			proc.cmd = exec.Command(
				"claude",
				"-p", proc.message,
				"--resume", proc.sessionId,
			)
			proc.cmd.Env = append(proc.cmd.Environ(),
				fmt.Sprintf("PENELOPE_WF_ID=%v", proc.workflowId),
				fmt.Sprintf("PENELOPE_SESSION_ID=%v", proc.sessionId),
			)
			proc.cmd.Stderr = &proc.stderr
			proc.cmd.Stdout = &proc.stdout

			if err := proc.cmd.Start(); err != nil {
				proc.logger.Error("error sending message to claude session",
					slog.Any("err", err),
					slog.String("stderr", proc.stderr.String()),
				)
				return err
			}
			proc.logger.Info("stdout",
				slog.String("stdout", proc.stdout.String()))
			return nil
		},
	}
	send.Flags().StringVar(&proc.sessionId, "session-id", "", "session id to resume")
	send.Flags().StringVarP(&proc.message, "message", "m", "", "message to send")
	_ = send.MarkFlagRequired("session-id")
	_ = send.MarkFlagRequired("message")

	sig := &cobra.Command{
		Use: "signal",
		RunE: func(cmd *cobra.Command, args []string) error {
			proc.logger.Info("penelope signal",
				slog.String("sessionId", proc.sessionId),
				slog.String("workflowId", proc.workflowId),
			)
			return nil
		},
	}
	sig.Flags().StringVar(&proc.sessionId, "session-id", "", "session id to signal")
	sig.Flags().StringVar(&proc.workflowId, "workflow-id", "", "workflow id")
	_ = sig.MarkFlagRequired("session-id")
	_ = sig.MarkFlagRequired("workflow-id")

	root.AddCommand(start, send, sig)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
