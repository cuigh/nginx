此镜像是为了解决 SPA 前端项目在 Docker 下的多环境部署问题，它可以在站点启动时把环境依赖参数自动注入到项目中。

## 快速上手

example 目录下包含了一个完整的示例项目，对一个已有的 SPA 项目来说，你可以按照下面的步骤做改造。

### 1. 添加配置文件

在`src/config`目录下添加各环境下的配置文件，并按`app.{profile}.json`规范命名，如`app.production.json`

```json
{
    "API_BASE_URL": "https://{your-api-site-domain}"
}
```

### 2. 在项目中使用配置

由于启动器会把环境变量注入到 `window.config` 全局变量中，所以项目中使用环境变量的方式需要做一些调整。

```javascript
const baseUrl = window.config?.API_BASE_URL || import.meta.env.VITE_API_BASE_URL;
```

如果你使用了 TypeScript，你可能还需要扩展`Window`对象的定义以通过编译，在 config 目录中新建一个 config.d.ts 文件，并输入如下内容

```typescript
interface Window {
    config?: {
        API_BASE_URL: string
    }
}
```

现在你只需要在`.env.development`或`.env.local`中添加开发环境的配置参数即可，其它环境下的配置你应该把它们挪到`config`目录中。

### 3. 添加 Dockerfile 文件

假设项目输出目录为 dist，把下面三行代码复制到 Dockerfile 文件中即可构建可支持多环境的通用镜像。

```dockerfile
FROM cuigh/nginx
COPY dist .
COPY src/config config/
```

### 4. 启动站点

启动容器时通过`profile`环境变量激活指定环境的配置参数。

```shell
docker run -e PROFILE=prd --name test -it -p 8080:80 test
```

### 5. 高级功能

1. 如果你喜欢，你可以使用 yaml 格式的配置文件。
2. 如果你喜欢，你依然可以使用 .env.{profile} 文件来存放配置，从而继续享受美好的旧时光。

需要注意的是，无论使用哪种配置文件，在构建 Docker 镜像时一定要把配置文件包含进去。