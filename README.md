#### rtmp2flv

![](./docs/images/rtmp2flvad.png)

##### 项目功能：

1. rtmp转httpflv播放
2. rtmp视频录像，录像文件为flv格式

##### 运行说明：

1. 下载[程序文件](https://github.com/hkmadao/rtmp2flv/releases)，解压   
2. 执行程序文件：window下执行rtmp2flv.exe，linux下执行rtmp2flv   
3. 浏览器访问程序服务地址：http://[server_ip]:9090/rtmp2flv/#/
4. 在网页配置摄像头推送的code、rtmpAuthCode等信息(如：rtmp://127.0.0.1:1935/camera/9527,则code为：camera,rtmpAuthCode为：9527)  
5. 等待摄像头连接，观看视频      

> 注意：
>
> ​	程序目前支持h264视频编码、aac音频编码，若不能正常播放，关掉摄像头推送的音频再尝试

##### 目录结构：

```
--rtmp2flv #linux执行文件
--rtmp2flv.exe #window执行文件
  --static #程序的网页文件夹
  --conf #配置文件文件夹
    --conf.yml #配置文件
  --db #sqlite3 #数据库文件夹
    --rtmp2flv.db #sqlite3数据库文件（存放摄像头推送的code、rtmpAuthCode等信息）
  --output #程序输出文件夹
    --live #保存摄像头录像的文件夹，录像格式为flv
    --log #程序输出的日志文件夹
```

##### 配置说明：

```
server:
    user:
        name: admin #网页登录用户名
        password: admin #网页登录密码
    rtmp:
        port: 1935 #程序的http端口
    httpflv:
        port: 9090
    fileflv:
        save: true #是否保存录像文件
        path: ./output/live #录像文件夹
    log:
        path: ./output/log #日志文件夹        
```

##### 开发说明：

程序分为服务器和页面，服务端采用golang开发，前端采用react+materia-ui，完成后编译页面文件放入服务器的static文件夹

###### 服务器开发说明：

1. 安装golang，MinGW(sqlite3模块的使用到的cgo,window下开发需要使用，window下可选择安装MinGW)
2. 获取[服务器源码](https://github.com/hkmadao/rtmp2flv.git)
3. 进入项目目录
4. go build开发

###### 页面开发说明：

1. 安装node
2. 下载[页面源码](https://github.com/hkmadao/rtmp2flv-web.git)
3. 进入项目目录
4. npm install
5. npm run start