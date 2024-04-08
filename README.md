# Box

## Todo

configx: 可配置多种配置源，如文件、etcd等。定义基本的配置文件格式，定义加载mysql、redis等组件的位置，实现bizConfig的追踪

launcher：配合configx，实现对各种组件和配置的自动加载，简化入口代码。优雅停机功能。

wire：配合开源的wire，实现各种组件的依赖注入和加载，减少手动组装的代码

components：以单例模式启动各种组件，如MySQL，Redis等。增加对各种组件的open telemetry监控

cachex：多级缓存。本地缓存可配置使用不同的开源库，中心缓存使用redis，为list和map等结构定义AutoFetch方法。增加防止缓存击穿、本地和中心缓存的强一致性等。

responsor: 各种情况下返回gin或grpc请求的方法合集。

middlewares：一些公共的中间件

## proto生成器

gin

grpc

openapi

validate

## 其他生成器

gorm-gen.yaml
