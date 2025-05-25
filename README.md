# WebDAV Docker部署

1. 根据你的需要修改`main.go`
    > [!WARNING]
    > 不要在`main.go`中设置端口号和共享路径
    ```go
    // ...
    // 设置WebDAV访问用户名
    username := "username"
    // 设置WebDAV访问密码
    password := "password"
    // ...
   ```
2. 将本项目复制到服务器
3. 使用下面的命令来创建镜像，镜像名称随意
   ```bash
   docker build -t dav <本项目在服务器中的位置>
   ```
4. 使用下面的命令创建容器
   ```bash
   sudo docker run -d \
   --restart always \
   -v <主机WebDAV映射路径>:/app/dir
   -p <WebDAV在主机上访问的端口号>:3000
   -name dav dav
   ```