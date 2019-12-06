# Change Log

## v2.1.1

### Bug Fixes

+ 新路由中间件省缺路径，仅指定方法
+ 由于 endless 无法在 Windows 下运行，单独给 linux 写平滑切换

### Features

+ rest 适配 new router
+ auth 适配 new router
+ 增加 Makefile
+ 增加热更新工具 `sola-hot`

#### Next Version

+ v2.2.x
    + 反向代理功能完善

### Breaking Changes

+ router & x/router 使用的 ctx key 已更改，请使用包提供的标准方法操作数据！

## v2.1.0-beta (2019-12-03)

### Bug Fixes

无

### Features

+ 新路由中间件取代旧路由中间件（进行中）

#### Next Version

+ v2.1.1
    + auth 适配 new router
    + rest 适配 new router
    + ...
    + router & x/router 使用的 ctx key 更改
+ v2.2.x
    + 反向代理功能完善

### Breaking Changes

+ router -> x/router
+ xrouter -> router

## v2.0.0-alpha (xxxx-xx-xx)

### Bug Fixes

+ 修复了路由组件匹配上的问题
+ 修复了 auth 中间件 Cookie path 的问题导致的认证失败
+ Option 为 nil 的问题

### Features

+ 合并了 v1 中的一些中间件
+ 新增了中间件： cors、swagger、logger、rest、graphql、(ws)WebSocket、新router
+ 增加独立的 box 仓库存放扩展中间件
+ Context Shadow
+ Use 允许多中间件

### Breaking Changes

+ 改变了中间件的定义，并移到 sola 包中。
