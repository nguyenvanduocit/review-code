package cli

import (
	"converter/pkg/gemini"
	"converter/pkg/git"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code using AI",
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.Diff()
		if err != nil {
			return err
		}

		fmt.Println(diff)

		model, closeGemini, err := gemini.NewGeminiClient(cmd.Context())
		if err != nil {
			return err
		}
		defer closeGemini()

		contents := []genai.Part{
			genai.Text("Generate code review for the following diff:"),
			genai.Text(diff),
		}

		aiResponse, err := model.GenerateContent(cmd.Context(), contents...)
		if err != nil {
			return err
		}

		outputText := strings.TrimSpace(string(aiResponse.Candidates[0].Content.Parts[0].(genai.Text)))
		cmd.Println(outputText)
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
