## [urlooker](https://github.com/710leo/urlooker)
监控web服务可用性及访问质量，采用go语言编写，易于安装和二次开发    

## Feature
- 返回状态码检测
- 页面响应时间检测
- 页面关键词匹配检测
- 自定义Header
- GET、POST、PUT访问
- 自定义POST BODY
- 检测结果支持推送 nightingale、open-falcon

## Architecture
![Architecture](img/urlooker_arch.png)

## 常见问题
- [wiki手册](https://github.com/710leo/urlooker/wiki)
- [常见问题](https://github.com/710leo/urlooker/wiki/FAQ)
- 初始用户名密码：admin/password

## Install
#### docker 安装

```bash
git clone https://github.com/710leo/urlooker.git
cd urlooker
docker build .
docker volume create urlooker-vol
docker run -p 1984:1984 -d --name urlooker --mount source=urlooker-vol,target=/var/lib/mysql --restart=always [CONTAINER ID]
```

#### 源码安装
```bash
# 插入数据库表
wget https://raw.githubusercontent.com/710leo/urlooker/master/sql/schema.sql
mysql -h 127.0.0.1 -u root -p < schema.sql

# 从github上拉取对应的代码
git clone -b develop https://github.com/404notfoundf/urlooker.git

# 将web.yml配置文件中mysql部分，修改用户名密码
 mysql:
    addr: "username:password@tcp(127.0.0.1:3306)/urlooker?charset=utf8&&loc=Asia%2FShanghai"
    idle: 10
    max: 20
# 如果需要添加不同时间的访问, 首先需要在原来已经存在的url的基础上进行配置, 在agent.yml添加配置即可（url支持数组）
  url_interval:
    - url: ["http://www.abc.com", "http://www.a.com"]
      interval: 60
    - url: ["http://www.abd.com"]
      interval: 90

# 编译及运行
./control build web
./control build agent

./control start web
./control start agent

# 注：如果添加不同时间访问，需要重新启动agent
./control restart agent
```
打开浏览器访问 http://127.0.0.1:1984 即可