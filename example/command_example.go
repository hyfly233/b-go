package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 定义要执行的 Python 脚本和参数
	cmd := exec.Command("python3", "script.py", "arg1", "arg2")

	// 获取命令的输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// 打印输出
	fmt.Println(string(output))
}
