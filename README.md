# Deep-Packet-Inspection<br>

Windows下配置gopacket环境:<br>
https://studygolang.com/articles/12116<br>
解决collect2.exe: error ld returned 1 exit status错误所需要的.a文件已置于根目录<br>

编译提示缺少类库，直接 `go get 类库url` 即可

功能测试:<br>
抓取、解析数据报<br>
创建、发送TCP或UDP报文<br>
创建、发送DNS报文<br>
DNS服务器（for captive portal）<br>
HTTP/2推送服务器<br>
