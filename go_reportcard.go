// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type GoReportCardGrade struct{}

func NewGoReportCardGrade(_ *Config) (Messenger, error) {
	return &GoReportCardGrade{}, nil
}

func (m GoReportCardGrade) Message() (string, error) {
	var buf bytes.Buffer
	var bufErr strings.Builder

	cmd := exec.Command("goreportcard-cli", "-j")
	cmd.Stdout = &buf
	cmd.Stderr = &bufErr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running goreportcard-cli (%w)", err)
	}

	errOut := bufErr.String()
	if errOut != "" {
		return "", fmt.Errorf("goreportcard-cli produced messages in STDERR")
	}

	reportCard := &goReportCard{}
	if err := json.Unmarshal(buf.Bytes(), reportCard); err != nil {
		return "", fmt.Errorf("unmarshalling goreportcard-cli output (%w)", err)
	}

	return reportCard.GradeFromPercentage, nil
}

type goReportCard struct {
	GradeFromPercentage string `json:"GradeFromPercentage"`
}
