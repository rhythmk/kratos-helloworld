# 1、创建student表

```sql
CREATE TABLE `student` (
  `id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;


```
同时修改 修改config mysql 的连接信息配置。
# 2、进入gorm
