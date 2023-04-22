## A Light-weight TCP Server

这是一个轻量级的 TCP 服务器框架，他只有 1 千多行代码，是参考 [zinx 开源项目]("https://github.com/aceld/zinx")实现的玩具框架，但已经包含了一个 TCP 服务框架的核心功能，如并发连接、多路由、读写分离、业务工作池等，能够很容易地搭建高性能并发连接服务器。

1. 并发连接：一个 go 程监听连接请求，将每个连接分配给新 go 程处理；
2. 多路由：根据 message（二进制消息格式为：len|id|data）的 id 号索引路由，找到对应的 handler 函数；
3. 读写分离：将每个连接的 socket 读写分到两个 go 程中处理，由管道传输 data；
4. 业务工作池：分配设定数量的 worker 来调用 handler 完成业务处理，以避免产生同时存在大量业务 go 程时的高资源压力。
5. 核心参数配置：以 json 格式配置在服务器的 config/config.json 目录中，见 zdemo/Server/config/config.json

## Simple IM Server/Client

基于该框架实现的一个单机版的实时聊天服务器以及相应的 terminal 客户端，支持群聊和群内成员的私聊两种模式。User.go 将连接和用户绑定，Room.go 维护聊天室的群聊和里面用户的增删查改等，World.go 维护所有的聊天室的增删查。各种 Router.go 是具体的业务 Handler

当前非常简单地在 map 中把用户和群聊的 name 和对象进行绑定，因此不能重名，例如修改名字时不能和群聊中的人重名，切换群聊时如果和新群中的用户重名则会无法切换（需要更名）。这些缺点应该很容易通过将 IP 和用户名绑定进行解决，只是这样查找用户名的时候需要遍历一波，没有这样实现。

客户端支持的命令有

1. 群聊
   bc\nline1\nline2\n...\nlineN\neof
2. 私聊
   to [user name]\nline1\nline2\n...\nlineN\neof
3. 修改聊天名称[默认用户名称是 ip 地址，不能和群内成员重名]
   rename [your new name]
4. 查看当前群聊的所有用户
   whos
5. 创建新群
   newromm [room name] [capacity]
6. 切换群聊
   enter [room name]
7. 查看所有群聊状态，名称 当前人数 最大支持人数
   rooms
8. 查看当前所在群名
   curroom
