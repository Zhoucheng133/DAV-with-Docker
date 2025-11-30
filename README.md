# WebDAV Docker部署

## Docker Hub

```bash
sudo docker run -d \
--restart always \
-v <主机WebDAV映射路径>:/app/dir \
-e USERNAME=<WebDAV用户名> \
-e PASSWORD=<WebDAV密码> \
-p <WebDAV在主机上访问的端口号>:3000 \
zhouc1230/webdav:latest
```

## 手动部署

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