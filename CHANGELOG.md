# Change Log

## v2.1.2

### Bug Fixes

+ sola-hot 监听 *.yml
+ 使用 Handler/Middleware 互转等工具方法优化代码，防止 nil 错误
+ hot print log only in log mode

### Features

+ `(*router.Router).Bind(pattern string, h sola.Handler, ms ...sola.Middleware)`
+ Default Sola App (`DefaultApp`)
+ `sola` -> `Listen*` 支持多 App 参数，默认监听 `DefaultApp`
+ 增加 `sola.ListenKeep`/`sola.ListenKeepTLS`
+ Handler/Middleware 互转工具方法，Handler/Middleware Must 防止 nil 错误
+ C/H/M 语法糖
+ Pass Handler
+ 内部错误 & InternalError Handler

### Breaking Changes

+ v2.1.0 & v2.1.1 内部版本号变量误标记为 `2.0.0`，从 v2.1.2 开始该变量格式变更为 `v2.x.x`
+ 移除 `native.From`/`native.FromFunc`，使用 `sola.FromHandler`/`sola.FromHandlerFunc`
+ 移除 `(Handler).Adapter() func(http.ResponseWriter, *http.Request, error)`，适配模块采用 `(Handler).ToHandler`/`(Handler).ToHandlerFunc`

### Next Version

+ v2.1.3
    + 检查所有 `M()` 是否需要替换为 `Must()`
+ v2.2.x
    + context 核心化
    + 移除 x/router

## v2.1.1 (2019-12-11)

### Bug Fixes

+ 新路由中间件省缺路径，仅指定方法
+ 由于 endless 无法在 Windows 下运行，单独给 linux 写平滑切换

### Features

+ rest 适配 new router
+ auth 适配 new router
+ 增加 Makefile
+ 增加热更新工具 `sola-hot`
+ 增加动态模块加载扩展
+ Context 获取 Store (包括 Origin 的 Store)
+ 反向代理功能完善 - 负载均衡

#### Next Version

+ v2.2.x
    + context 核心化
    + 移除 x/router

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
