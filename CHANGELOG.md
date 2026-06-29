## [1.9.1](https://github.com/codepzj/shortvid-backend/compare/v1.9.0...v1.9.1) (2026-06-29)


### Bug Fixes

* 调整rustfs的sts临时凭证为最小权限 ([cf84afe](https://github.com/codepzj/shortvid-backend/commit/cf84afe3fdda8768517e37e039198a71fbc46cc2))

# [1.9.0](https://github.com/codepzj/shortvid-backend/compare/v1.8.0...v1.9.0) (2026-06-28)


### Features

* user-service修改http响应格式 ([6a072ce](https://github.com/codepzj/shortvid-backend/commit/6a072ce2c683f321541f69812237b4a264e2ba53))

# [1.8.0](https://github.com/codepzj/shortvid-backend/compare/v1.7.0...v1.8.0) (2026-06-27)


### Features

* 完善用户服务(升级kratos为v3) ([7aaa098](https://github.com/codepzj/shortvid-backend/commit/7aaa098fdeb04136ba2f48d301004cd4f0b654ca))
* 完善短视频主服务(升级kratos为v3) ([d696e83](https://github.com/codepzj/shortvid-backend/commit/d696e8381115d850f1c2fd9117648a9573867d69))
* 添加user微服务 ([a893a5f](https://github.com/codepzj/shortvid-backend/commit/a893a5f6c86b0de483ae81def933acac97241d96))

# [1.7.0](https://github.com/codepzj/shortvid-backend/compare/v1.6.0...v1.7.0) (2026-06-24)


### Features

* 添加日志 ([e1b0163](https://github.com/codepzj/shortvid-backend/commit/e1b0163348bca97819a1b643d76815e65133e863))

# [1.6.0](https://github.com/codepzj/shortvid-backend/compare/v1.5.2...v1.6.0) (2026-06-21)


### Features

* 添加health路由 ([3ba2040](https://github.com/codepzj/shortvid-backend/commit/3ba20405c630f73d004ab1d624a5802d7973b670))

## [1.5.2](https://github.com/codepzj/shortvid-backend/compare/v1.5.1...v1.5.2) (2026-06-21)


### Bug Fixes

* 使用cnb ([b6700c1](https://github.com/codepzj/shortvid-backend/commit/b6700c1dcbc33940574c0152763c220f32b4a29a))

## [1.5.1](https://github.com/codepzj/shortvid-backend/compare/v1.5.0...v1.5.1) (2026-06-21)


### Bug Fixes

* 修改deploy配置 ([4a5c9bd](https://github.com/codepzj/shortvid-backend/commit/4a5c9bdd1ad366f9c50ca684f8de52c5c8660051))

# [1.5.0](https://github.com/codepzj/shortvid-backend/compare/v1.4.0...v1.5.0) (2026-06-21)


### Features

* 手动选择tag ([dade294](https://github.com/codepzj/shortvid-backend/commit/dade29459deab25fa856d70396763ab45f552475))

# [1.4.0](https://github.com/codepzj/shortvid-backend/compare/v1.3.0...v1.4.0) (2026-06-21)


### Features

* release后触发docker打包 ([af0a740](https://github.com/codepzj/shortvid-backend/commit/af0a74068dd0ad5a987b84ad83622201e24b22f8))

# [1.3.0](https://github.com/codepzj/shortvid-backend/compare/v1.2.0...v1.3.0) (2026-06-21)


### Features

* 触发tag ([4d97edd](https://github.com/codepzj/shortvid-backend/commit/4d97edd4d67c8dea1bf5f699a0cc7b86f5a52c45))

# [1.2.0](https://github.com/codepzj/shortvid-backend/compare/v1.1.0...v1.2.0) (2026-06-21)


### Features

* 添加k8s部署方式 ([b9cd4a2](https://github.com/codepzj/shortvid-backend/commit/b9cd4a2d8f4383b634efbc2b49ae0395bf60a066))

# [1.1.0](https://github.com/codepzj/shortvid-backend/compare/v1.0.1...v1.1.0) (2026-06-18)


### Features

* 合并yml ([cc12bf4](https://github.com/codepzj/shortvid-backend/commit/cc12bf46e80e6b34d4b67f4c0f9c67ec25199eca))

## [1.0.1](https://github.com/codepzj/shortvid-backend/compare/v1.0.0...v1.0.1) (2026-06-18)


### Bug Fixes

* ci写权限 ([7228e0c](https://github.com/codepzj/shortvid-backend/commit/7228e0c8259414672d57b32f5fb93b24a7457e4d))

# 1.0.0 (2026-06-18)


### Bug Fixes

* firebase从infra移除 ([6d8dd31](https://github.com/codepzj/shortvid-backend/commit/6d8dd316116fc7d0c9e453be0c5b79af8b8cd477))
* logger使用helper ([ef73290](https://github.com/codepzj/shortvid-backend/commit/ef732909c0f65b8bdf736f35daef538f1f35db03))
* 修复ci问题 ([5e9a3e9](https://github.com/codepzj/shortvid-backend/commit/5e9a3e94e03ac94995c04d09f27888ba9a0c6b99))
* 修复更新用户信息的错误 ([2b6345a](https://github.com/codepzj/shortvid-backend/commit/2b6345afb7dea87dea1c31f44b1299637e8728e3))
* 修复获取不到github info的问题 ([72bbc0b](https://github.com/codepzj/shortvid-backend/commit/72bbc0bdd18ed26cdb5b26e39e2de5e17c288766))
* 删除不必要的函数 ([7ec209e](https://github.com/codepzj/shortvid-backend/commit/7ec209ea0bae1a849646e0909384dcd5058be141))
* 去除file service ([9e8a5f4](https://github.com/codepzj/shortvid-backend/commit/9e8a5f4fe8db3b7ef5f00184a94bce9e73eecd97))
* 引入账号体系, 变更sql和mvc层代码 ([bc2e5f6](https://github.com/codepzj/shortvid-backend/commit/bc2e5f6708b5e5247fca9f4b9babb67b0b44a6f3))
* 调整路由 ([703af14](https://github.com/codepzj/shortvid-backend/commit/703af14dd2c8f933d2849f37a516d46e4a27f6b4))


### Features

* http&grpc服务器添加参数验证中间件 ([2f410cf](https://github.com/codepzj/shortvid-backend/commit/2f410cf45cd3a42c57a5b8362c16838e5daebac9))
* logger转为helper ([a7a0c06](https://github.com/codepzj/shortvid-backend/commit/a7a0c06493c937ae1ab66c51aacc901879ec9ed4))
* mysql添加最大打开连接数和最大空闲连接数的配置 ([225eedb](https://github.com/codepzj/shortvid-backend/commit/225eedbacc92b53d7b758dec18d9adf49a07e320))
* user和account创建强关联 ([20dd8b3](https://github.com/codepzj/shortvid-backend/commit/20dd8b3193585eb30944e5e722632dd04a77893a))
* 上传服务使用s3通用配置 ([fc2d345](https://github.com/codepzj/shortvid-backend/commit/fc2d345205ce85febca1f50cbbfdd18ad1f6aee6))
* 使用s3客户端 ([dbd5b1d](https://github.com/codepzj/shortvid-backend/commit/dbd5b1d0e6e3fef1fe79b70614b8e5cee261790c))
* 删除自定义logger, 新增公参中间件 ([8d1996b](https://github.com/codepzj/shortvid-backend/commit/8d1996be0e25987e3c232c6a0fb4da3b0879ab0d))
* 完成github登录 ([2d64698](https://github.com/codepzj/shortvid-backend/commit/2d646986a18df40d7a69c2a716466c599559e609))
* 引入账号体系 ([632864e](https://github.com/codepzj/shortvid-backend/commit/632864e97ce4a5ef7fabfd0431dd0709ccfe3eac))
* 接入mysql ([868409a](https://github.com/codepzj/shortvid-backend/commit/868409ad656de2489979bfa60e6a385c3a5afbb3))
* 注入minio ([66385e7](https://github.com/codepzj/shortvid-backend/commit/66385e7bdf435c3e3cd84dcc9579c21f8eaee122))
* 添加ci环境 ([1cb29cc](https://github.com/codepzj/shortvid-backend/commit/1cb29cc7685d59f68632fd91e482262b0f42c6af))
* 添加firebase登录接口 ([892c3d9](https://github.com/codepzj/shortvid-backend/commit/892c3d90c581058d9ef51418652604e521a2c481))
* 添加github oauth2获取用户信息逻辑 ([4018169](https://github.com/codepzj/shortvid-backend/commit/401816940676b48deff319183a2ac8eeaaa4d010))
* 添加github登录接口 ([e187c63](https://github.com/codepzj/shortvid-backend/commit/e187c63c9386a5914bcb2a19b428ff09d46a3c20))
* 添加jwt鉴权和路由白名单机制 ([2395145](https://github.com/codepzj/shortvid-backend/commit/2395145c1644adbc160d20bcf57080758136f775))
* 添加session,cookie机制 ([a66ac82](https://github.com/codepzj/shortvid-backend/commit/a66ac82f277d3909b23c958cf30357ff004053d9))
* 添加shortvid-service的docker部署 ([d43d340](https://github.com/codepzj/shortvid-backend/commit/d43d34004a5829e1cad818913f862d826afd8a37))
* 添加video建表语句 ([9753a58](https://github.com/codepzj/shortvid-backend/commit/9753a58ea19233c370e3e966323e6114cfce8926))
* 添加zap日志 ([97f906c](https://github.com/codepzj/shortvid-backend/commit/97f906c64d04f88b37300354f6039dd4f65fde50))
* 添加公共参数中间件 ([91e4c02](https://github.com/codepzj/shortvid-backend/commit/91e4c0246f4d27026e783205aaabd279586b26a5))
* 添加文件上传服务 ([b89895e](https://github.com/codepzj/shortvid-backend/commit/b89895e548c361fa34f1f236f449971d65a0fb1c))
* 添加校验 ([779513c](https://github.com/codepzj/shortvid-backend/commit/779513c192dd9c246d02f885abbf0d5efb90de86))
* 添加用户UID生成和登录信息更新功能 ([d1a466a](https://github.com/codepzj/shortvid-backend/commit/d1a466aa7df6b4f931628be262f2e8a3114e2224))
* 添加用户服务 ([4e3e630](https://github.com/codepzj/shortvid-backend/commit/4e3e630adee985296983d7c079483e7e53aa1fca))
* 添加错误码和自定义http响应格式 ([2c48332](https://github.com/codepzj/shortvid-backend/commit/2c48332142deac659a5fcb268f908fdbb826e3f8))
* 用户响应添加accessToken和refreshToken ([0cc47b4](https://github.com/codepzj/shortvid-backend/commit/0cc47b449d99c644a8d7575e165633b2d1e4bd47))
* 重构HTTP服务器的鉴权中间件 ([7222597](https://github.com/codepzj/shortvid-backend/commit/7222597869dbf81d2e2426eb33ccac1dadde55d5))
