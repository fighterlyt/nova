package stack

/*
Stack 堆栈操作，将递归转为循环，并且记录路径
参数:
*	init   	      T          	初始化参数
*	process	      func(T) []T	递归计算
*	get    	      func(T) V  	从数据获取路径,可以为空，为空时不返回路径
* initCapacity  int         用户缓存的slice 初始容量
返回值:
*	[]V    	      []V        	路径，按照T出现的顺序从早到晚排列
*/
func Stack[T, V any](init T, process func(T) []T, get func(T) V, initCapacity int) []V {
	var (
		stackItem = make([]T, 0, initCapacity)
		path      []V
	)

	stackItem = append(stackItem, init)

	for len(stackItem) > 0 {
		next := stackItem[len(stackItem)-1]
		stackItem = stackItem[:len(stackItem)-1]

		if get != nil {
			path = append(path, get(next))
		}

		data := process(next)

		stackItem = append(stackItem, data...)
	}

	if get != nil {
		return path
	}

	return nil
}
