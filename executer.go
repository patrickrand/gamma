package gamma

import (
	"encoding/json"
	"os/exec"
)

type Executer interface {
	Execute(cmd string, args ...string) (code int, message string)
}

type ShellExecuter struct {
	defaultErrorCode int
}

func NewShellExecuter(defaultErrorCode int) *ShellExecuter {
	return &ShellExecuter{defaultErrorCode}
}

func (sh ShellExecuter) Execute(cmd string, args ...string) (code int, message string) {
	data, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return sh.defaultErrorCode, "gamma.ShellExecuter failed to execute the command: " + err.Error()
	}

	output := struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{}

	if err := json.Unmarshal(data, &output); err != nil {
		return sh.defaultErrorCode, "gamma.ShellExecuter failed to unmarshal the command's output: " + err.Error()
	}

	return output.Code, output.Message
}
