package commands

import (
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for your shell.
Examples:
  source <(csync completion bash)
  csync completion zsh > ~/.zshrc
		`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				logger.Logger.Warn().Msg("Unsupported shell. Please use bash, zsh, fish, or powershell.")
			}
		},
	}
}
