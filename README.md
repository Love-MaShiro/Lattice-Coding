创建单仓库go模块骨架

# 测试 API 服务
curl http://localhost:8080/api/v1/providers

# 测试 MySQL 连接
mysql -h localhost -P 3306 -u admin -ppassword --protocol=tcp -e "SELECT 1 as test"

# 测试 Redis 连接
redis-cli -h localhost -p 6379 ping

# 中止测试项目
docker stop lattice-api && docker rm lattice-api && docker rmi $(docker images -q 'lattice-coding*' 2>/dev/null) 2>/dev/null; cd /home/ubuntu/Lattice-Coding && docker build -t lattice-api:latest .