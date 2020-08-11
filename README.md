# WIP
无证施工中...
![image](res/workgolang.jpg)
# nactor
基于nats中间件开发的actor框架。

# 设计目的
1. 一般的游戏服务器不同于传统app，会有大量的实体对象之间的来回通信。通常的做法是游戏进程自己定制管理各对象的通信，非常复杂且难以扩展。如果基于actor框架，就可以在框架层解决了通信问题，提高了开发效率。  
但是go没有完善的actor框架,之前使用star最多的protoactor，不过等几年也没发布正式版本。   
2. 本着学习和减少以来的目的，将于来的ganat+protoactor框架合并，借助nats中间件完成此框架。  
3. 基于TDD原则开发并实践，提供代码质量

# 进度
- [x] actor: 主体
- [x] actor: 消息收发,tell/req,订阅额外主题
- [x] actor: watch机制,kill,死亡通知
- [ ] actor: 健康检查
- [x] service: 主体
- [x] service: rpc封装
- [x] service: 服务发现
- [ ] app:进程启动，服务管理，配置管理
- [ ] fuction:gateway
- [ ] fuction:robot
- [ ] fuction:room direct,route
- [ ] fuction:db redis/mongo
- [ ] fuction:自动化属性系统
- [ ] example:statefull
- [ ] example:stateless
- [ ] example:chat
