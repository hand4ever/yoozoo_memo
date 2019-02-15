## docker pratice 学习

### 1. 几个 docker 相关解释

**`docker run -it --rm ubuntu:16.04 bash`**

> -i: 交互式操作
>
> -t: 终端
>
> --rm 容器退出后随之将其删除
>
> ubuntu:16.04 镜像:tag
>
> bash 命令

**`docker system df`**

> 查看镜像、容器、数据卷所占用的空间

**`docker image ls -f dangling=true`**

> 显示虚悬镜像，同名pull或者build后被挤下为<none>的景象
> 可以用 `docker image prune`命令删除 dangling image

**`docker image ls -q`**

> 列出所有镜像 id

**分层存储**

> **Docker镜像**
>
> 操作系统分为*内核*和*用户空间*，对于 Linux 而言，内核启动后，会挂载 `root 文件系统`为其提供用户空间支持。Docker 镜像，就相当于是一个 `root 文件系统`，除了提供容器运行时所需的程序、库、资源、配置等文件外，还包含一些为运行时准备的配置参数（如 匿名卷、环境变量、用户等）。镜像不包含任何动态数据，其内容在构建后也不会被改变。
>
> **分层存储**
>
> 利用 Union FS 技术，Docker 设计为分层存储的架构。
>
> 镜像在构建时，会一层层构建，前一层是后一层的基础。每一层构建完就不会再发生改变，后一层的任何改变只发生在自己这一层。在构建镜像时，尽量只包含该层需要添加的东西，任何额外的东西应该在该层构建结束前清理掉。
>
> Union FS 是有最大层数限制的，比如 `AUFS`，目前是最大 `127` 层

### 2. 使用 docker commit 命令，手动给旧的镜像添加了新的一层，形成新的镜像

> docker commit 可以在容器被入侵后后保存现场，定制镜像应该使用 Dockerfile 来完成，*慎用 docker commit*，因为这些操作是黑箱操作，除了制作的人知晓执行了什么命令，如何生成，别人无从知晓。

**step 1. `docker run --name webserver -d -p 80:80 nginx`**

> 启动一个容器，命令为 webserver，并且映射了 80 端口，然后用 http://localhost/ 来访问

**step 2. `docker exec -it webserver bash`**

> 进入容器，准备修改内容
> `echo '<h1>Hello, nginx; Hi, Docker!</h1>' > /usr/share/nginx/html/index.html`

**step 3. 保存镜像，即将容器的存储层保存下来，成为镜像**

```sh
docker commit \
--author "lpan@yoozoo.com" \
--message "修改了默认网页内容" \
webserver \
nginx:v2
```

> 用`docker history nginx:v2`来查看镜像内的历史记录

**step 4. 运行新的定制好的镜像**
```sh
docker run --name web2 -d -p 81:80 nginx:v2
```

### 3. 使用 Dockerfile 定制镜像

> Dockerfile 本身是一个文本文件，其内容包含了一条条指令，每一条指令构建一层。

**step 1. 新建目录，然后新建 Dockerfile 文件，内容如下：**

```dockerfile
FROM nginx
RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
```

> `必须` **FROM** 指定基础镜像
>
> *服务类镜像*，如 `nginx`、`redis`、`mysql`、`httpd`、`php`、`tomcat`等；
>
> *构建、运行各种语言应用类镜像*，如`node`、`python`、`golang`等
>
> *还有更基础的操作系统镜像*，如 `ubuntu`、`debian`、`centos`等
>
> *空白镜像* ：`scratch`；
>
> ​	这意味着不以任何镜像为基础，接下来所写的指令将作为镜像的第一层开始存在。
>
> ​	     如 `swarm`、 `coreos/etcd`

> **RUN** 执行命令
>
> *shell 格式*： `RUN <命令>`
>
> *exec 格式*  `RUN ["可执行文件", "参数1", "参数2"]`

> - Dockerfile 支持 Shell 那样的行尾添加 `\` 来换行
>
> - 行首 `#` 用来注释
>
> - 每一层执行完要执行清理工作

**step 2 在 Dockerfile 文件所在目录执行命令**

```sh
docker build -t nginx:v3 .
#输出结果如下
Sending build context to Docker daemon  2.048kB
Step 1/2 : FROM nginx
 ---> f09fe80eb0e7
Step 2/2 : RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
 ---> Running in 190ad3f8a2a4
Removing intermediate container 190ad3f8a2a4
 ---> 15ad186d2e09
Successfully built 15ad186d2e09
Successfully tagged nginx:v3
```

> 命令格式 
>
> ​	`docker build [选项] <上下文路径>`
>
> 工作原理
>
> ​	Docker 在运行时分为 Docker 引擎和客户端工具。
>
> ​	  Docker 引擎提供了一组 REST API，被称为 Docker Remote API，而如 `docker`命，令这样的客户端工具，则是通过这组 API 与 Docker 引擎交互，来完成各种功能。因此，我们在本机上执行各种 `docker`功能，实际上是**使用远程调用形式在服务端（Docker 引擎）完成**。
>
>