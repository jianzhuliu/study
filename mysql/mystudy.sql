
/*多行注释*/
#单行注释
--单行注释

--创建用户
create user 'canal'@'%' identified by '123456';
grant select, replication slave,replication client on *.* to 'canal'@'%'; 
alter user 'canal'@'%' identified by 'canal';



--创建数据库
show create database mysql;
create database mystudy default character set utf8mb4 collate utf8mb4_0900_ai_ci;
use mystudy;

--创建表
DROP TABLE IF EXISTS test1;
CREATE TABLE test1 (
 id INT NOT NULL COMMENT '编号',
 name VARCHAR(20) NOT NULL COMMENT '姓名',
 sex TINYINT NOT NULL COMMENT '性别,1男,2女',
 email VARCHAR(50)
);

--准备数据
DROP PROCEDURE IF EXISTS proc1;
DELIMITER $
CREATE PROCEDURE proc1()
BEGIN
	DECLARE i INT DEFAULT 1;
	START TRANSACTION;
	WHILE i <= 2000000 DO
		INSERT INTO test1 (id, name, sex, email) VALUES (i,concat('javacode',i),if(mod(i,2),1,2),concat('javacode',i,'@163.com'));
		SET i = i + 1;
		if i%10000=0 THEN
			COMMIT;
			START TRANSACTION;
		END IF;
	END WHILE;
	COMMIT;
END $
DELIMITER ;

--调用存储过程
CALL proc1();

--校验
SELECT count(1) FROM test1 limit 1;

--非索引下查询 
select * from test1 where id = 1 limit 1;

select * from test1 where name = 'javacode1';

select * from test1 a where a.email = 'javacode1000085@163.com';

--创建索引，然后查询
create index idx1 on test1 (id);
select * from test1 where id = 1 limit 1;

create unique index idx2 on test1(name);
select * from test1 where name = 'javacode1';

create index idx3 on test1 (email(15));
select * from test1 a where a.email = 'javacode1000085@163.com';

--删除索引
 drop index idx1 on test1;
