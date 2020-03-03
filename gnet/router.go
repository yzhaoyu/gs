package gnet

import "github.com/yzhaoyu/gs/giface"

// 实现router时，先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写即可
type BaseRouter struct{}

// 这里之所以BaseRouter的方法都为空是因为有的Router不希望有PreHandle、PostHandle
// 所以Router全部继承BaseRouter的好处就是不需要实现PreHandle、PostHandle

// 在处理conn业务之前的钩子方法
func (br *BaseRouter) PreHandle(request giface.IRequest) {}

// 在处理conn业务的主方法
func (br *BaseRouter) Handle(request giface.IRequest) {}

// 在处理conn业务之后的钩子方法
func (br *BaseRouter) PostHandle(request giface.IRequest) {}
