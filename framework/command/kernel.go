package command

import "github.com/gohade/hade/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(initEnvCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initAppCommand())
	root.AddCommand(initBuildCommand())
	root.AddCommand(initDevCommand())
	root.AddCommand(initProviderCommand())
	root.AddCommand(initCmdCommand())
	root.AddCommand(initSwaggerCommand())
}