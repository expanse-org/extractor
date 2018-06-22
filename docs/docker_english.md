# Loopring Extractor Docker

The loopring development team provides the loopring/extractor image. The latest latest version is v1.5.0

## Run
get the latest docker image
``` 
docker pull loopring/extractor
```
create log&config dir
```bash
mkdir your_log_path your_config_path
```
config extractor.toml [参考](https://loopring.github.io/extractor/EXTRACTOR_DOCUMENT_CN)<br>

before deployment, perform telnet tests according to the ports related to the configuration files mysql, redis, kafka, and zk to ensure that these dependencies can be accessed normally.

mount the log dir and config dir, and run
```bash
docker run --name extractor -idt -v your_log_path:/opt/loopring/extractor/logs -v your_config_path:/opt/loopring/extractor/config loopring/extractor:latest --config=/opt/loopring/extractor/config/extractor.toml /bin/bash
```

## History version

| version         | desc         |
|--------------|------------|
| v1.5.0| the first release version|
