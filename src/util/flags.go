package util

// ParseFlag 用于解析flag，并且返回是否继续执行程序流程，false为退出
// 会在启动web服务器之前运行
// 如果error不为nil时，也会打印错误信息，然后退出
func ParseFlag() (bool, error) {
	return true, nil
}
