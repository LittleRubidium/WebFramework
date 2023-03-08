package console

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/command"
)

func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		//定义根命令的关键字
		Use: "hade",
		//简短介绍
		Short: "hade命令",
		//根命令的详细介绍
		Long: "hade 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令 ",
		//根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		//不需要出现 cobra默认的completion命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	//为根Command设置服务容器
	rootCmd.SetContainer(container)
	//绑定框架命令
	command.AddKernelCommands(rootCmd)
	//绑定业务命令

	//执行RootCommand
	return rootCmd.Execute()
}
