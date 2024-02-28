package cli

import (
	"converter/pkg/convert"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.PersistentFlags().String("obsidian-path", "v", "Absolute path to the obsidian vault")
	convertCmd.PersistentFlags().String("hugo-path", "v", "Absolute path to the hugo site directory")

	if err := viper.BindPFlags(convertCmd.PersistentFlags()); err != nil {
		panic(err)
	}
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert obsidian notes to hugo",
	RunE: func(cmd *cobra.Command, args []string) error {
		obsidianPath := viper.GetString("obsidian-path")
		if obsidianPath == "" {
			return cmd.Help()
		}

		hugoPath := viper.GetString("hugo-path")
		if hugoPath == "" {
			return cmd.Help()
		}

		return convert.ObsidianToHugo(obsidianPath, hugoPath)
	},
}
