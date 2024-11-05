# GinChat
仿微信聊天界面的聊天平台

# GinChat README

## 1. 项目概述
本项目是一个仿微信聊天界面的聊天平台，使用Gin框架进行开发，提供了一系列的聊天相关功能接口。

## 2. 项目结构
- `docs`目录可能包含了用于生成API文档的相关信息。
- `models`目录包含了与数据模型相关的代码，如处理音频上传、群聊连接等功能相关的模型代码。
- `service`目录包含了业务逻辑服务相关的代码，用于处理如获取首页、用户操作、消息发送等功能。

## 3. 路由设置
### 3.1 总体路由设置
`Router`函数是整个项目的路由设置核心函数，它创建了一个`gin.Engine`实例，并进行了一系列的配置和路由注册。
1. **API文档和静态资源配置**
   - 通过`docs.SwaggerInfo.BasePath = ""`和相关路由注册，使得可以通过`/swagger/*any`路径访问Swagger生成的API文档。
   - 设置静态资源目录为`/asset`，并加载`views/**/*`的HTML模板文件。
2. **页面路由设置**
   - 提供了多个页面相关的路由，如`/`、`/index`、`/toRegister`、`/toChat`等，分别对应不同的页面操作，通过`service`模块中的相关函数来处理页面请求。
3. **用户相关路由**
   - 包括获取用户列表（`/user/getUserList`）、创建用户（`/user/createUser`）、删除用户（`/user/deleteUser`）、更新用户（`/user/updateUser`）、登录（`/user/login`）等操作，通过`service`模块中的相关函数来实现具体功能。
4. **消息相关路由**
   - 包括发送消息（如`/user/sendMsg`、`/user/sendUserMsg`）、搜索好友（`/user/searchFriends`）、加好友（`/user/addFriend`）、加群（`/user/addGroup`）等操作，同样通过`service`模块中的相关函数来处理。
   - 还包括上传附件（`/attach/upload`）、音频上传处理（`/user/audioUploadHandler`）等功能相关的路由。
5. **群聊相关路由**
   - 包括创建群聊（`/user/createCommunity`）、加载群聊列表（`/contact/loadCommunity`）、保存消息（`/contact/saveMessage`）、查询最近消息记录（`/contact/getRecentMessages`）等操作，通过`service`模块中的相关函数来实现群聊相关功能。
   - 同时设置了群聊连接的路由（`/groupChat`）和消息广播的路由（`/sendGroupMessage`），分别由`models`模块中的相关函数来处理。

## 4. 如何运行
1. 确保已经安装了项目所需的依赖，如`gin`、`swaggo`等相关库。
2. 配置好项目所需的环境，可能包括数据库连接（从代码中可以推测可能使用了MongoDB，但未明确）等相关配置。
3. 运行项目的入口文件（如果有的话，从提供的代码中未明确），项目将启动并监听相应的端口（未明确具体端口，但可以推测如果是本地开发，可能是默认的`gin`端口或者通过配置文件指定的端口）。

## 5. 接口文档
项目中集成了Swagger文档，可以通过访问`/swagger/*any`路径查看接口文档，了解各个接口的详细参数和返回值信息。