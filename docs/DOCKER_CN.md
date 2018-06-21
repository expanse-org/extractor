# Loopring Extractor Docker 中文文档

## 简介
loopring开发团队提供loopring/extractor镜像,最新版本是v1.5.0。<br>

## 部署
* 获取docker镜像
```bash
docker pull loopring/extractor
```
* 创建log&config目录
```bash
mkdir your_log_path your_config_path
```
* 配置extractor.toml文件，[参考](https://loopring.github.io/extractor/EXTRACTOR_DOCUMENT_CN)
* telnet测试mysql,redis,zk,kafka,ethereum rpc相关端口能否连接

## 运行
```bash
docker run --name extractor -idt -v your_log_path:/opt/loopring/extractor/log -v your_config_path:/opt/loopring/extractor/config loopring/extractor:latest --config=/opt/loopring/extractor/config/extractor.toml /bin/bash
```

## 历史版本
| 版本号         | 描述         |
|--------------|------------|
| v1.5.0| release初始版本|
