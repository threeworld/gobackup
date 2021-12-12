package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	path        = flag.String("path", "", "Traverse the directory under the path compression path, delete after compression.")
	thread      = flag.Int("t", 2, "Set number of concurrent coroutines.")
	spaceRegexp = regexp.MustCompile(`[\s]+`)
)

func main() {

	flag.Parse()
	if flag.NFlag() < 1 {
		log.Fatal("Please enter path parameters")
	}
	dirs, err := FindDir(*path)

	var wg sync.WaitGroup
	ch := make(chan struct{}, *thread)
	if err == nil {
		err = chPwd(*path)
		fmt.Println(os.Getwd())
		if err == nil {
			for _, dir := range dirs {
				ch <- struct{}{}
				wg.Add(1)
				go func(dir string) {
					output, err := compressTar(dir+".tar.gz", dir)
					if err != nil {
						fmt.Println("tar fail", output, err)
						panic(err)
					}
					output, err = rmDir(dir)
					if err != nil {
						fmt.Println("tar fail", output, err)
						panic(err)
					}
					fmt.Println("tar and rm success", path, dir)
					<-ch
					defer wg.Done()
				}(dir)
			}
		}
	}
	wg.Wait()
}

func FindDir(path string) ([]string, error) {
	var dirs []string
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return dirs, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}

func Exec(command string, args ...string) (output string, err error) {
	//去除空格并且转为切片
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	cmdArgs := []string{}
	//判断命令是否存在
	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}
	//如果传入的命令带参数
	if len(commands) > 1 {
		cmdArgs = commands[1:]
	}
	//拼接传入的参数
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, args...)
	}

	cmd := exec.Command(fullCommand, cmdArgs...)
	cmd.Env = os.Environ()
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	stdOut, err := cmd.Output()
	fmt.Println(fullCommand, " ", strings.Join(cmdArgs, " "))
	if err != nil {
		err = errors.New(stdErr.String())
		return
	}
	output = strings.Trim(string(stdOut), "\n")
	return
}

// compressTar  压缩文件
//  @param fileName 压缩后的文件名
//  @param day   压缩的目录
func compressTar(fileName, dir string) (string, error) {
	tarGz := "tar -zcvf"
	output, err := Exec(tarGz, fileName, dir)
	return output, err
}

// rmDir       删除压缩的目录
//  @param day  压缩的目录
func rmDir(dir string) (string, error) {
	rm := "rm -rf"
	output, err := Exec(rm, dir)
	return output, err
}

// chPwd 切换工作目录
//  @param pwd 需要切换到的目录
func chPwd(pwd string) error {

	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		if pwd == execDir {
			return nil
		} else if err := os.Chdir(pwd); err != nil {
			return err
		}
	}
	return err
}
