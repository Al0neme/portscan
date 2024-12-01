# portscan

- 简介

通过http请求探测端口是否开放，避免全端口开放的情况，判断的依据是是否超时以及响应是否是EOF异常，所以可能存在其他异常导致的判断失误，如果遇到大批不正常的端口结果，请手动验证端口是否真的开放

- 参数

```bash
┌─┐┌─┐┬─┐┌┬┐┌─┐┌─┐┌─┐┌┐┌
├─┘│ │├┬┘ │ └─┐│  ├─┤│││
┴  └─┘┴└─ ┴ └─┘└─┘┴ ┴┘└┘                                   
Author: Al0neme

Usage of ./portscan:
  -i string
        host addr
  -p string
        ports, example:80,443 or 8000-9000 or -
  -s int
        connect timeout, default 3s (default 3)
  -t int
        number of thread, default 10 (default 10)
```

- 使用
  
```bash
git clone https://github.com/Al0neme/portscan.git
cd portscan
go build .
./portscan -i 127.0.0.1 -p - -t 100
```
