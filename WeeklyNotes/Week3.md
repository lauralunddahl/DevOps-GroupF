## Week 3

### DB abstraction layer in your ITU-MiniTwit
We have chosen GORM as our Object-relational mapping, and will be using MySQL as the DBMS.
Reasoning for choosing the GORM library and the MySQL database management service is that ...

Commands:
To install MySQL: sudo apt install mysql-server

Creating and filling the DB locally for now (until we have it up and running remotely):
```
sudo mysql
create database minitwit;
use minitwit;
create table follower (who_id int, whom_id int);
create table message (message_id int primary key auto_increment, author_id int not null, text varchar(250) not null, pub_date int, flagged int);
create table user (user_id int primary key auto_increment, username varchar(100) not null, email varchar(50) not null, pw_hash varchar(500) not null);
desc 
GRANT ALL PRIVILEGES ON * . * TO 'groupf'@'localhost';
```

You must also get these:
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm

Setup:
```
export GOPATH='[path to repo]'
export GOROOT=''
```

After having done this, you will have to go get all of the packages again that we have gotten so far. 
