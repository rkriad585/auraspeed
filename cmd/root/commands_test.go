package root

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRootCmd(t *testing.T) {
	output, err := executeCommand(rootCmd)
	if err != nil {
		t.Errorf("unexpected error executing root command: %v", err)
	}
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestRootCmdHelp(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestConfigViewCmd(t *testing.T) {
	cmd := newConfigViewCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigViewCmdWithSection(t *testing.T) {
	cmd := newConfigViewCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"speedtest"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigSetCmdInvalidArgs(t *testing.T) {
	cmd := newConfigSetCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"onlyonearg"})

	err := cmd.Execute()
	if err == nil {
		t.Errorf("expected error for invalid args, got nil")
	}
}

func TestConfigSetCmdTooManyArgs(t *testing.T) {
	cmd := newConfigSetCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"key", "value", "extra"})

	err := cmd.Execute()
	if err == nil {
		t.Errorf("expected error for too many args, got nil")
	}
}

func TestHistoryCmd(t *testing.T) {
	cmd := NewHistoryCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHistoryCmdWithLimit(t *testing.T) {
	cmd := NewHistoryCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--limit", "5"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHistoryCmdWithSaveFlag(t *testing.T) {
	cmd := NewHistoryCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--save"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHistoryCmdWithClearFlag(t *testing.T) {
	cmd := NewHistoryCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--clear"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSpeedtestCmdHelp(t *testing.T) {
	cmd := NewSpeedtestCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestSpeedtestCmdFlagParsing(t *testing.T) {
	cmd := NewSpeedtestCommand()

	serverIDFlag := cmd.Flags().Lookup("server-id")
	if serverIDFlag == nil {
		t.Errorf("expected --server-id flag to be defined")
	}

	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag == nil {
		t.Errorf("expected --json flag to be defined")
	}
}

func TestSpeedtestCmdInvalidServerID(t *testing.T) {
	cmd := NewSpeedtestCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--server-id", "notanumber"})

	err := cmd.Execute()
	if err == nil {
		t.Errorf("expected error for invalid server ID, got nil")
	}
}

func TestInfoCmdHelp(t *testing.T) {
	cmd := NewInfoCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestNetworkCmdHelp(t *testing.T) {
	cmd := NewNetworkCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestNetworkPingCmdHelp(t *testing.T) {
	cmd := NewNetworkCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"ping", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestNetworkTracerouteCmdHelp(t *testing.T) {
	cmd := NewNetworkCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"traceroute", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestNetworkDnsCmdHelp(t *testing.T) {
	cmd := NewNetworkCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"dns", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestConfigCmdHelp(t *testing.T) {
	cmd := NewConfigCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestConfigViewCmdHelp(t *testing.T) {
	cmd := NewConfigCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"view", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestConfigSetCmdHelp(t *testing.T) {
	cmd := NewConfigCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"set", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}

func TestConfigResetCmdHelp(t *testing.T) {
	cmd := NewConfigCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"reset", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("expected help output, got empty string")
	}
}
