package cobra

import "github.com/gohade/hade/framework"

//设置服务器容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

//获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
