package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func AssignWork(info []*Line) {
	numberCPU := 100
	lenInfo := len(info)
	everyGoWork := 3
	AssignWorkChan := make(chan []*Line, numberCPU)
	wg := sync.WaitGroup{}
	Yes := make(chan *Line, 100)
	Bad := make(chan *Line, 100)
	yesYes := make([]*Line, 0, 3000)
	badBad := make([]*Line, 0, 1500)
	warpDone := make(chan struct{})
	JnDu := make(chan int, 100)

	go WrapWork(Yes, Bad, &yesYes, &badBad, warpDone, JnDu)
	go PrintTiao(JnDu, lenInfo)

	for i := 0; i < numberCPU; i++ {
		wg.Add(1)
		go Work(AssignWorkChan, &wg, Yes, Bad)
	}

	start := 0
	end := 0
	for {
		end += everyGoWork
		if end >= lenInfo {
			end = lenInfo
			break
		}
		AssignWorkChan <- info[start:end]
		start = end
	}
	AssignWorkChan <- info[start:end]
	close(AssignWorkChan)

	wg.Wait()
	close(Bad)
	close(Yes)
	<-warpDone
	close(warpDone)
	close(JnDu)
	fmt.Println("")

	// 写入文件

	// 打开文件，使用 O_WRONLY|O_CREATE|O_TRUNC 标志来覆盖写入
	// 如果文件不存在，则创建文件；如果文件存在，则清空文件内容
	file, err := os.OpenFile("Success.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("Error opening file:" + err.Error())
		return
	}

	// 将字符串写入到缓冲写入器中
	for _, value := range yesYes[:len(yesYes)-1] {
		_, err = io.WriteString(file, (value.Id + "," + value.Name + "," + value.Personal + "," + value.Grid + "," + value.Address + "\n"))
	}
	_, err = io.WriteString(file, (yesYes[len(yesYes)-1].Id + "," + yesYes[len(yesYes)-1].Name + "," + yesYes[len(yesYes)-1].Personal + "," + yesYes[len(yesYes)-1].Grid + "," + yesYes[len(yesYes)-1].Address))

	if err != nil {
		panic("Error writing to buffer:" + err.Error())
		return
	}

	fmt.Println("Successfully wrote to success file.")
	file.Close() // 确保文件在函数返回时被关闭

	// 打开文件，使用 O_WRONLY|O_CREATE|O_TRUNC 标志来覆盖写入
	// 如果文件不存在，则创建文件；如果文件存在，则清空文件内容
	file, err = os.OpenFile("Unsuccess.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("Error opening file:" + err.Error())
		return
	}

	// 将字符串写入到缓冲写入器中
	for _, value := range badBad[:len(badBad)-1] {
		_, err = io.WriteString(file, (value.Id + "," + value.Name + "," + value.Personal + "," + value.Grid + "," + value.Address + "\n"))
	}
	_, err = io.WriteString(file, (badBad[len(badBad)-1].Id + "," + badBad[len(badBad)-1].Name + "," + badBad[len(badBad)-1].Personal + "," + badBad[len(badBad)-1].Grid + "," + badBad[len(badBad)-1].Address))
	if err != nil {
		panic("Error writing to buffer:" + err.Error())
		return
	}

	fmt.Println("Successfully wrote to unsuccess file.")
	file.Close() // 确保文件在函数返回时被关闭
}

func WrapWork(yes chan *Line, bad chan *Line, yesYes *[]*Line, badBad *[]*Line, done chan struct{}, JinDu chan int) {
	bool1 := true
	bool2 := true
	yesLine := &Line{}
	badLine := &Line{}
	allCount := 0

	for {
		select {
		case yesLine, bool1 = <-yes:
			if yesLine != nil {
				//fmt.Println("\033[32;4m", yesLine, "\033[0m")
				*yesYes = append(*yesYes, yesLine)
				allCount++
				if allCount%10 == 0 {
					fmt.Println(" 成功示例: \033[32;4m", yesLine, "\033[0m")
					JinDu <- allCount
				}
			}
		case badLine, bool2 = <-bad:
			if badLine != nil {
				//fmt.Println("\033[31;4m", badLine, "\033[0m")
				*badBad = append(*badBad, badLine)
				allCount++
				if allCount%10 == 0 {
					fmt.Println("\033[31;4m 失败示例: ", badLine, "\033[0m")
					JinDu <- allCount
				}
			}
		}
		if !bool1 && !bool2 {
			break
		}
	}
	done <- struct{}{}
}

func Work(infos chan []*Line, wg *sync.WaitGroup, yes chan *Line, bad chan *Line) {
	client := http.Client{
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Second * 3,
		},
	}
	for info := range infos {
	FOR2:
		for _, web := range info {

			url := web.Address
			if len(url) <= 4 {
				bad <- web
				continue FOR2
			}
			if url[0:4] != "http" {
				bad <- web
				continue FOR2
			}

			resp, err := client.Get(url)

			if err != nil {
				//fmt.Println("请求失败:", err)
				if resp != nil {
					err1 := resp.Body.Close()
					if err1 != nil {
						panic("http连接关闭失败: " + err.Error())
					}
				}
				bad <- web
				continue FOR2
			}

			// 检查状态码
			if resp.StatusCode == http.StatusOK {
				yes <- web
				//fmt.Println("请求成功，状态码为200")
			} else {
				bad <- web
				//fmt.Printf("请求失败，状态码为: %d\n", resp.StatusCode)
			}

			func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					panic("http连接关闭失败: " + err.Error())
				}
			}(resp.Body) // 确保关闭响应体
		}
	}
	wg.Done()
}

func PrintTiao(control chan int, all int) {
	for {
		select {
		case JinDu, ok := <-control:
			if !ok {
				println("\033[1F\033[s\033[1000F" + "\033[36m\033[2K" + "                           全部完成!" + "\033[0m\033[u")
				return
			}
			Percent := float64(JinDu) / float64(all)
			Percent = Percent * 100
			println("\033[1F\033[s\033[1000F" + "\033[36m\033[2K" + "                           已完成: " + strconv.FormatFloat(Percent, 'f', 3, 32) + "% \033[0m\033[u")
		default:
			println("\033[1F\033[s\033[1000F\033[1E\033[2K\033[u")
		}
	}
}
