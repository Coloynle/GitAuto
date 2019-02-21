package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type GitAuto struct {
	WorkDir string
	Branch  string
}

// 主函数 命令行接收参数 第一个参数为当前目录下的项目名 第二个参数为分支名
// 默认和所有项目同级目录 运行更新所有项目 如果从某个项目路径运行 则只更新本项目
func main() {
	var Git GitAuto
	var argsLen int = len(os.Args)
	if 1 == argsLen {
		RunPath, _ := os.Getwd()
		FilePath := Git.getCurrentPath()
		RunPath = RunPath + "\\"
		if RunPath != FilePath {
			Git.WorkDir = RunPath
			Git.updateProject()
		} else {
			Git.AllProject()
		}
	} else if 2 == argsLen {
		Git.WorkDir = Git.getCurrentPath() + os.Args[1]
		Git.updateProject()
	} else if 3 == argsLen {
		Git.WorkDir = Git.getCurrentPath() + os.Args[1]
		Git.Branch = os.Args[2]
		Git.updateProject()
	}
}

// 更新所有Git项目
func (G *GitAuto) AllProject() {
	var rootPath string = G.getCurrentPath()
	var dir []string = G.getDir()
	for _, d := range dir {
		G.WorkDir = rootPath + d
		G.updateProject()
	}
}

// 更新Git项目
func (G *GitAuto) updateProject() {
	fmt.Print("*****************************************************************\n")
	fmt.Print("                  " + G.WorkDir + "\n")
	fmt.Print("*****************************************************************\n")
	var NowBranch string
	if "" == G.Branch {
		NowBranch = G.Ioutil()
	} else {
		NowBranch = G.Branch
	}
	if "0" == NowBranch {
		fmt.Print("This is not a GIT project\n")
		return
	}
	G.gitReset()
	G.gitCheckout("master")
	G.gitFetch()
	G.gitRebase("origin/master")
	if "master" != NowBranch {
		G.gitClean()
		G.gitCheckout(NowBranch)
		G.gitRebase("origin/" + NowBranch)
	}
}

// 获取当前完整路径
func (G *GitAuto) getCurrentPath() string {
	s, _ := exec.LookPath(os.Args[0])
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

// 获取当前目录下所有非隐藏文件夹
func (G *GitAuto) getDir() []string {
	var dir []string
	files, _ := ioutil.ReadDir(G.getCurrentPath())
	for _, f := range files {
		if f.IsDir() {
			str := f.Name()
			matched, err := regexp.MatchString("^\\.\\S*", str)
			if err == nil && !matched {
				dir = append(dir, f.Name())
			}
		}
	}
	return dir
}

// 获取Git项目当前分支
func (G *GitAuto) Ioutil() string {
	var name string = G.WorkDir + "/.git/HEAD"
	if contents, err := ioutil.ReadFile(name); err == nil {
		// 因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		result = strings.Replace(result, "ref: refs/heads/", "", 1)
		return result
	}
	return "0"
}

// 重置项目
func (G *GitAuto) gitReset() {
	commandName := "git"
	params := []string{"reset", "--hard"}
	G.execCommand(commandName, params)
}

// 获得最新代码
func (G *GitAuto) gitFetch() {
	commandName := "git"
	params := []string{"fetch"}
	G.execCommand(commandName, params)
}

// 更新本地分支到最新分支
func (G *GitAuto) gitRebase(branch string) {
	commandName := "git"
	params := []string{"rebase", branch}
	G.execCommand(commandName, params)
}

// 切换分支
func (G *GitAuto) gitCheckout(branch string) {
	commandName := "git"
	params := []string{"checkout", branch}
	G.execCommand(commandName, params)
}

// 清除多余文件
func (G *GitAuto) gitClean() {
	commandName := "git"
	params := []string{"clean", "-df"}
	G.execCommand(commandName, params)
}

// CD
func (G *GitAuto) cd(path string) {
	commandName := "cd"
	params := []string{path}
	G.execCommand(commandName, params)
}

// 调用命令
func (G *GitAuto) execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	cmd.Dir = G.WorkDir

	// 显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	// 实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	return true
}
