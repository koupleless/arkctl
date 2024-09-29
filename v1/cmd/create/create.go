package create

import (
	"bytes"
	"fmt"
	"github.com/koupleless/arkctl/v1/cmd/root"
	"github.com/spf13/cobra"
	"os/exec"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Convert to module auto-configuration",
	Long: `create 命令用于自动模块化。
它会执行以下操作:
1. 修改 application.properties 文件
2. 创建 bootstrap.properties 文件
3. 修改 pom.xml 文件

使用方法:
  arkctl create -p <项目路径> -a <应用名称>

示例:
  arkctl create -p /path/to/project -a myapp`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("projectPath")
		applicationName, _ := cmd.Flags().GetString("applicationName")

		jarPath := "D:\\koupleless-ext-appModuleAutomator\\target\\demo-0.0.1-SNAPSHOT.jar"
		runJavaProgram(jarPath, projectPath, applicationName)
	},
}

// runJavaProgram 运行 Java 程序
func runJavaProgram(jarPath, projectPath, applicationName string) {
	// 创建一个命令来运行 Java 程序
	javaCmd := exec.Command("java", "-jar", jarPath)

	// 创建一个管道，用于传递输入到 Java 程序
	stdinPipe, err := javaCmd.StdinPipe()
	if err != nil {
		fmt.Printf("创建标准输入管道失败: %s\n", err)
		return
	}

	// 将标准输出和标准错误重定向到缓冲区
	var outBuffer bytes.Buffer
	javaCmd.Stdout = &outBuffer
	javaCmd.Stderr = &outBuffer

	// 启动 Java 程序
	if err := javaCmd.Start(); err != nil {
		fmt.Printf("启动 Java 程序出错: %s\n", err)
		return
	}

	// 写入项目路径到 Java 程序的标准输入
	if _, err := stdinPipe.Write([]byte(projectPath + "\n")); err != nil {
		fmt.Printf("写入项目路径出错: %s\n", err)
		return
	}

	// 写入应用名称到 Java 程序的标准输入
	if _, err := stdinPipe.Write([]byte(applicationName + "\n")); err != nil {
		fmt.Printf("写入应用名称出错: %s\n", err)
		return
	}

	// 关闭标准输入管道
	if err := stdinPipe.Close(); err != nil {
		fmt.Printf("关闭标准输入管道出错: %s\n", err)
		return
	}

	// 等待 Java 程序执行完成
	if err := javaCmd.Wait(); err != nil {
		fmt.Printf("Java 程序运行出错: %s\n", err)
		fmt.Printf("错误输出: %s\n", outBuffer.String())
		return
	}

	fmt.Printf("Java 程序输出:\n%s\n", outBuffer.String())
}

func init() {
	root.RootCmd.AddCommand(createCmd)

	// 定义命令行参数
	createCmd.Flags().StringP("projectPath", "p", "", "项目路径 (必填)")
	createCmd.Flags().StringP("applicationName", "a", "", "应用名称 (必填)")

	// 标记必填参数
	createCmd.MarkFlagRequired("projectPath")
	createCmd.MarkFlagRequired("applicationName")
}
