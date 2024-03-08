package git

import (
	"github.com/nguyenvanduocit/executils"
	"os"
	"strings"
)

func Diff() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	out := strings.Builder{}
	if err := executils.Run("git",
		executils.WithDir(workingDir),
		executils.WithArgs("diff", "--cached", "--unified=0"),
		executils.WithStdOut(&out),
	); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
