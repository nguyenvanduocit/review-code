package git

import (
	"github.com/nguyenvanduocit/executils"
	"os"
)

func Add() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	return executils.Run("git",
		executils.WithDir(workingDir),
		executils.WithArgs("add", "."),
	)
}
