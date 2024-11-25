# 使用官方MySQL 5.7镜像作为基础镜像
FROM mysql:5.7

# 复制自定义配置文件到容器中
COPY ./mysql.cnf /etc/mysql/conf.d/

# 设置配置文件权限
RUN chmod 644 /etc/mysql/conf.d/mysql.cnf 