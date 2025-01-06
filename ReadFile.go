package main

import (
	"bufio"
	"os"
	"strings"
)

type Line struct {
	Id       string
	Name     string
	Personal string
	Grid     string
	Address  string
}

func AssignmentLine() func(param string) *Line {
	a := &Line{}
	i := -1
	return func(param string) *Line {
		i++
		switch i {
		case 0:
			a.Id = param
		case 1:
			a.Name = param
		case 2:
			a.Personal = param
		case 3:
			a.Grid = param
		case 4:
			a.Address = param
		case 5:
			return a
		default:
			return a
		}
		return a
	}
}

func ReadFile(path string) (res []*Line) {
	res = make([]*Line, 0, 3600)
	file, err := os.Open(path)
	if err != nil {
		panic("Error opening file:" + err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("Error closing file:" + err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)

	// 逐行读取文件
FOR1:
	for scanner.Scan() {
		line := scanner.Text()             // 获取当前行的内容
		fields := strings.Split(line, ",") // 按逗号分隔每一行

		backup := make([][]string, 0, 5)
		backup = append(backup, fields)

		count := 0
		count += len(fields)
		bools := false
	FOR2:
		for {
			if len(fields) < 5 && count < 5 {
				scanner.Scan()
				line = scanner.Text()
				fields = strings.Split(line, ",")
				backup = append(backup, fields)
				count += len(fields)
				if bools {
					count--
				}
				bools = true
			} else {
				count = 0
				bools = false
				break FOR2
			}
		}

		if len(backup) == 1 {
			res = append(res, &Line{
				Id:       fields[0],
				Name:     fields[1],
				Personal: fields[2],
				Grid:     fields[3],
				Address:  fields[4],
			})
			continue FOR1
		}

		for i := 1; i < len(backup); i++ {
			backup[i][0] = backup[i-1][len(backup[i-1])-1] + backup[i][0]
		}
		backup = append(backup, []string{backup[len(backup)-1][len(backup[len(backup)-1])-1]})

		test := AssignmentLine()
		for _, val := range backup {
			for _, val1 := range val[:len(val)-1] {
				test(val1)
			}
		}
		test(backup[len(backup)-1][len(backup[len(backup)-1])-1])
		res = append(res, test(""))
	}

	// 检查是否发生错误
	if err := scanner.Err(); err != nil {
		panic("Error reading file:" + err.Error())
	}

	return res
}
