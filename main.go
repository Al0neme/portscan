package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func banner() {
	page := `
┌─┐┌─┐┬─┐┌┬┐┌─┐┌─┐┌─┐┌┐┌
├─┘│ │├┬┘ │ └─┐│  ├─┤│││
┴  └─┘┴└─ ┴ └─┘└─┘┴ ┴┘└┘                                   
Author: Al0neme
`
	fmt.Println(page)
}

func randomuseragent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 10; Pixel 3 XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/89.0",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := r.Intn(len(userAgents))
	return userAgents[randomIndex]
}

func checkPort(addr string, timeout int) bool {
	url := fmt.Sprintf("http://%s", addr)

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", randomuseragent())
	req.Header.Set("Origin", url)
	req.Header.Set("Referer", url)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		//判断是否是EOF，如果是则未开放端口
		if strings.Contains(err.Error(), "EOF") {
			return false
		}
		//判断是否是超时异常，如果是则未开放端口
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return false
		}

		return true
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	return true
}

func getresult(addr string) string {
	outresult := fmt.Sprintf("%s open", addr)
	fmt.Println(outresult)
	return outresult
}

func saveresult(result string) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	if _, err := file.WriteString(result + "\n"); err != nil {
		return
	}
}

func syncrun(targets []string, timeout int, concurrencylimit int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrencylimit)
	for _, addr := range targets {
		wg.Add(1)
		sem <- struct{}{}
		go func(a string, t int) {
			defer wg.Done()
			isOpen := checkPort(addr, timeout)
			if isOpen {
				saveresult(getresult(addr))
			}
			<-sem
		}(addr, timeout)
	}
	wg.Wait()
}

func parsePort(portarg string) []int {
	var ports []int

	// 如果输入是 "-"
	if portarg == "-" {
		for i := 1; i <= 65535; i++ {
			ports = append(ports, i)
		}
		return ports
	}

	// 分割输入
	parts := strings.Split(portarg, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			// 处理范围
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(rangeParts[0])
			end, err2 := strconv.Atoi(rangeParts[1])
			if err1 == nil && err2 == nil {
				for i := start; i <= end; i++ {
					ports = append(ports, i)
				}
			}
		} else {
			// 处理单个数字
			if num, err := strconv.Atoi(part); err == nil {
				ports = append(ports, num)
			}
		}
	}
	return ports
}

func main() {
	banner()

	host := flag.String("i", "", "host addr")
	port := flag.String("p", "", "ports, example:80,443 or 8000-9000 or -")
	thread := flag.Int("t", 10, "number of thread, default 10")
	timeout := flag.Int("s", 3, "connect timeout, default 3s")

	flag.Parse()
	if *host == "" {
		fmt.Println("usage: portscan -i 127.0.0.1 -p 22,80 -t 5")
		return
	}
	if *port == "" {
		fmt.Println("usage: portscan -i 127.0.0.1 -p 22,80 -t 5")
		return
	}
	config := fmt.Sprintf("thread: %d\ntimeout: %d\n", *thread, *timeout)
	fmt.Println(config)
	targets := make([]string, 0)
	ports := parsePort(*port)
	for _, port := range ports {
		targets = append(targets, fmt.Sprintf("%s:%d", *host, port))
	}

	syncrun(targets, *timeout, *thread)

	fmt.Println("all task done, the live target saved to result.txt")
}
