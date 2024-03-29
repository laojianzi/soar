# 重写规则

[toc]

## dml2select
* **Description**:将数据库更新请求转换为只读查询请求，便于执行EXPLAIN

* **Original**:

```sql
DELETE FROM film WHERE length > 100
```

* **Suggest**:

```sql
select * from film where length > 100
```
## reg2select
* **Description**:使用正则的方式将数据库更新请求转换为只读查询请求，便于执行EXPLAIN

* **Original**:

```sql
DELETE FROM film WHERE length > 100
```

* **Suggest**:

```sql
select * from film where length > 100
```
## star2columns
* **Description**:为SELECT *补全表的列信息

* **Original**:

```sql
SELECT * FROM film
```

* **Suggest**:

```sql
select film.film_id, film.title from film
```
## insertcolumns
* **Description**:为INSERT补全表的列信息

* **Original**:

```sql
insert into film values(1,2,3,4,5)
```

* **Suggest**:

```sql
insert into film(film_id, title, description, release_year, language_id) values (1, 2, 3, 4, 5)
```
## having
* **Description**:将查询的 HAVING 子句改写为 WHERE 中的查询条件

* **Original**:

```sql
SELECT state, COUNT(*) FROM Drivers GROUP BY state HAVING state IN ('GA', 'TX') ORDER BY state
```

* **Suggest**:

```sql
select state, COUNT(*) from Drivers where state in ('GA', 'TX') group by state order by state asc
```
## orderbynull
* **Description**:如果 GROUP BY 语句不指定 ORDER BY 条件会导致无谓的排序产生，如果不需要排序建议添加 ORDER BY NULL

* **Original**:

```sql
SELECT sum(col1) FROM tbl GROUP BY col
```

* **Suggest**:

```sql
select sum(col1) from tbl group by col order by null
```
## unionall
* **Description**:可以接受重复的时间，使用 UNION ALL 替代 UNION 以提高查询效率

* **Original**:

```sql
select country_id from city union select country_id from country
```

* **Suggest**:

```sql
select country_id from city union all select country_id from country
```
## or2in
* **Description**:将同一列不同条件的 OR 查询转写为 IN 查询

* **Original**:

```sql
select country_id from city where col1 = 1 or (col2 = 1 or col2 = 2 ) or col1 = 3;
```

* **Suggest**:

```sql
select country_id from city where (col2 in (1, 2)) or col1 in (1, 3);
```
## dmlorderby
* **Description**:删除 DML 更新操作中无意义的 ORDER BY

* **Original**:

```sql
DELETE FROM tbl WHERE col1=1 ORDER BY col
```

* **Suggest**:

```sql
delete from tbl where col1 = 1
```
## distinctstar
* **Description**:DISTINCT *对有主键的表没有意义，可以将DISTINCT删掉

* **Original**:

```sql
SELECT DISTINCT * FROM film;
```

* **Suggest**:

```sql
SELECT * FROM film
```
## standard
* **Description**:SQL标准化，如：关键字转换为小写

* **Original**:

```sql
SELECT sum(col1) FROM tbl GROUP BY 1;
```

* **Suggest**:

```sql
select sum(col1) from tbl group by 1
```
## mergealter
* **Description**:合并同一张表的多条ALTER语句

* **Original**:

```sql
ALTER TABLE t2 DROP COLUMN c;ALTER TABLE t2 DROP COLUMN d;
```

* **Suggest**:

```sql
ALTER TABLE t2 DROP COLUMN c, DROP COLUMN d;
```
## alwaystrue
* **Description**:删除无用的恒真判断条件

* **Original**:

```sql
SELECT count(col) FROM tbl where 'a'= 'a' or ('b' = 'b' and a = 'b');
```

* **Suggest**:

```sql
select count(col) from tbl where (a = 'b');
```
## countstar
* **Description**:不建议使用COUNT(col)或COUNT(常量)，建议改写为COUNT(*)

* **Original**:

```sql
SELECT count(col) FROM tbl GROUP BY 1;
```

* **Suggest**:

```sql
SELECT count(*) FROM tbl GROUP BY 1;
```
## innodb
* **Description**:建表时建议使用InnoDB引擎，非 InnoDB 引擎表自动转 InnoDB

* **Original**:

```sql
CREATE TABLE t1(id bigint(20) NOT NULL AUTO_INCREMENT);
```

* **Suggest**:

```sql
create table t1 (
	id bigint(20) not null auto_increment
) ENGINE=InnoDB;
```
## autoincrement
* **Description**:将autoincrement初始化为1

* **Original**:

```sql
CREATE TABLE t1(id bigint(20) NOT NULL AUTO_INCREMENT) ENGINE=InnoDB AUTO_INCREMENT=123802;
```

* **Suggest**:

```sql
create table t1(id bigint(20) not null auto_increment) ENGINE=InnoDB auto_increment=1;
```
## intwidth
* **Description**:整型数据类型修改默认显示宽度

* **Original**:

```sql
create table t1 (id int(20) not null auto_increment) ENGINE=InnoDB;
```

* **Suggest**:

```sql
create table t1 (id int(10) not null auto_increment) ENGINE=InnoDB;
```
## truncate
* **Description**:不带 WHERE 条件的 DELETE 操作建议修改为 TRUNCATE

* **Original**:

```sql
DELETE FROM tbl
```

* **Suggest**:

```sql
truncate table tbl
```
## rmparenthesis
* **Description**:去除没有意义的括号

* **Original**:

```sql
select col from table where (col = 1);
```

* **Suggest**:

```sql
select col from table where col = 1;
```
## delimiter
* **Description**:补全DELIMITER

* **Original**:

```sql
use sakila
```

* **Suggest**:

```sql
use sakila;
```
