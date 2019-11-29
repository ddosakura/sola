# Change Log

## v2.0.0-alpha (xxxx-xx-xx)

### Bug Fixes

+ 修复了路由组件匹配上的问题
+ 修复了 auth 中间件 Cookie path 的问题导致的认证失败

### Features

+ 合并了 v1 中的一些中间件
+ 新增了中间件： cors、swagger、logger、rest、graphql
+ 将一些 ctx 的扩展操作移到 x 包中（预发）
+ 增加独立的 box 仓库存放扩展中间件

### Breaking Changes

+ 改变了中间件的定义，并移到 sola 包中。
