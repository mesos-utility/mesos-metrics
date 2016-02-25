mesos集群监控脚本
================================

系统需求
--------------------------------
操作系统：Linux

主要逻辑
--------------------------------
获取mesos集群API接口数据，解析返回结果，将key组装成json后push到falcon-agent
接口解释请参照:
 * http://mesos.apache.org/documentation/latest/monitoring/
 * http://mesos.mydoc.io/

使用方法
--------------------------------
1. 根据实际部署情况，配置采集master或slave API接口;
 * master: http://127.0.0.1:5050/metrics/snapshot
 * slave:  http://127.0.0.1:5051/metrics/snapshot",

2. 测试： ./control build && ./control start
 * $GOPATH/bin/govendor init && $GOPATH/bin/govendor add +external && GO15VENDOREXPERIMENT=1 go build
