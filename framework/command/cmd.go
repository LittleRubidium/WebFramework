package command

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/util"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

var cmdCommand = &cobra.Command{
	Use: "command",
	Short: "控制台命令相关",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

//列出所有控制台命令
var cmdListCommand = &cobra.Command{
	Use: "list",
	Short: "列出所有控制台命令",
	RunE: func(c *cobra.Command, args []string) error{
		cmds := c.Root().Commands()
		ps := [][]string{}
		for _,cmd := range cmds {
			line := []string{cmd.Name(),cmd.Short}
			ps = append(ps,line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

//创建一个业务控制台命令
var cmdCreateCommand = &cobra.Command{
	Use: "new",
	Aliases: []string{"create","init"},
	Short: "创建一个控制台命令",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("开始创建控制台命令...")
		var name,folder string
		{
			prompt := &survey.Input{
				Message: "请输入控制台命令名称",
			}
			err := survey.AskOne(prompt,&name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同控制台命令):",
			}
			err := survey.AskOne(prompt,&folder)
			if err != nil {
				return err
			}
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder  := app.CommandFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		for _,subFolder := range subFolders {
			if subFolder == folder {
				fmt.Println("目录名称已存在")
				return nil
			}
		}

		//开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}

		//创建title这个模版方法
		funcs := template.FuncMap{"title":strings.Title}
		{
			file := filepath.Join(pFolder,folder, name+"go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			//使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTmpl))
			if err := t.Execute(f,name); err != nil {
				return errors.Cause(err)
			}
		}

		fmt.Println("创建新命令行工具成功，路径:",filepath.Join(pFolder,folder))
		fmt.Println("请记得开发完将命令行工具挂载到console/kernel.go")
		return nil
	},
}

// 命令行工具模版
var cmdTmpl string = `package {{.}}
import (
	"fmt"
	"github.com/gohade/hade/framework/cobra"
)
var {{.|title}}Command = &cobra.Command{
	Use:   "{{.}}",
	Short: "{{.}}",
	RunE: func(c *cobra.Command, args []string) error {
        container := c.GetContainer()
		fmt.Println(container)
		return nil
	},
}
`