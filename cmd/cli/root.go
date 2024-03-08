package cli

import (
	"codereview/pkg/gemini"
	"codereview/pkg/git"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "codereview",
	Short: "Review code using AI",
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.Diff()
		if err != nil {
			return err
		}

		model, closeGemini, err := gemini.NewGeminiClient(cmd.Context())
		if err != nil {
			return err
		}
		defer closeGemini()

		if len(diff) == 0 {
			fmt.Println("No changes to review, start by staging some changes with `git add`.")
			os.Exit(0)
		}

		contents := []genai.Part{
			genai.Text("You are a good code reviewer. Answer should be short, concise, straight forward. Do a code review for the following git diff, and provide feedback. ONLY response what need to be improve." + "\n" + "```" + "\n" + diff + "\n" + "```" + "\n" + "Feedback:" + "\n"),
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
