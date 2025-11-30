# WebDAV Docker

![License](https://img.shields.io/badge/License-MIT-dark_green)

This project is a WebDAV server deployed with Docker.  
Just install Docker and use the command to deploy!

Compatible with AMD64 & ARM64

[English](#english) | [中文](#中文)

## English

> [!NOTE]
> You need to replace the text inside the angle brackets (including the brackets themselves) in the following command.

### Build with Docker Hub

```bash
sudo docker run -d \
--restart always \
-v <host-path>:/app/dir \
-e USERNAME=<webdav-username> \
-e PASSWORD=<webdav-password> \
-p <host-port>:3000 \
--name dav \
zhouc1230/webdav:latest
```

### Manual Deployment
1. Dowload the project to your server
2. Build the image using the following command:
   ```bash
   docker build -t dav <project-path>
   ```
3. Use the following command to create a container
   ```bash
   sudo docker run -d \
   --restart always \
   -v <host-path>:/app/dir \
   -e USERNAME=<webdav-username> \
   -e PASSWORD=<webdav-password> \
   --name dav \
   -p <host-port>:3000 \
   dav
   ```

## 中文

> [!NOTE]
> 你需要修改下面命令中带有尖括号的内容（包括尖括号本身）

### Docker Hub

```bash
sudo docker run -d \
--restart always \
-v <主机WebDAV映射路径>:/app/dir \
-e USERNAME=<WebDAV用户名> \
-e PASSWORD=<WebDAV密码> \
-p <WebDAV在主机上访问的端口号>:3000 \
--name dav \
zhouc1230/webdav:latest
```

### 手动部署

1. 将本项目复制到服务器
2. 使用下面的命令来创建镜像，镜像名称随意
   ```bash
   docker build -t dav <本项目在服务器中的位置>
   ```
3. 使用下面的命令创建容器
   ```bash
   sudo docker run -d \
   --restart always \
   -v <主机WebDAV映射路径>:/app/dir \
   -e USERNAME=<WebDAV用户名> \
   -e PASSWORD=<WebDAV密码> \
   --name dav \
   -p <WebDAV在主机上访问的端口号>:3000 \
   dav
   ```