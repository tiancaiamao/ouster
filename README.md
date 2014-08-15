# 关于ouster
[![Build Status](https://travis-ci.org/tiancaiamao/ouster.svg?branch=master)](https://travis-ci.org/tiancaiamao/ouster)

ouster是一个用Go语言写的MMORPG服务端。

# 开发计划

以尽量小的代价做一个MMORPG原型，这是[最初的设计](http://zenlife.tk/ouster.slide)。目前已辞职全心投入此项目的开发。

# 下载和安装

配置好Go语言环境，获取服务端代码。

	go get github.com/tiancaiamao/ouster

客户端部分在[这里](http://pan.baidu.com/s/1sj4jIqT)，需要微软的视窗操作系统。
解压后，修改Data\Info目录下的ServerInfo.inf文件，将LoginServerAddress改成本机的IP。

将源文件ouster/config/config.go中的GameServerIP改为本机IP。

	go install github.com/tiancaiamao/ouster/gameserver
	go install github.com/tiancaiamao/ouster/loginserver

运行loginserver和gameserver，运行客户端。
