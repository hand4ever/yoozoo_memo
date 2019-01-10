## 背景
> gitlab 中使用 sentry，日志中间件使用

## 介绍
> [Sentry](https://sentry.io/welcome/) 是一个实时事件日志记录和汇集的平台。其专注于错误监控以及提取一切事后处理所需信息而不依赖于麻烦的用户反馈。它分为客户端和服务端，客户端(目前客户端有Python, PHP,C#, Ruby等多种语言)就嵌入在你的应用程序中间，程序出现异常就向服务端发送消息，服务端将消息记录到数据库中并提供一个web页方便查看。Sentry由python编写，源码开放，性能卓越，易于扩展，目前著名的用户有Disqus, Path, mozilla, Pinterest等。

## 目标

## 参考

 * [ 在PHP中使用Sentry](https://juejin.im/post/5a992115f265da239f06d0d7)
 * 

## 详细设计
### 安装