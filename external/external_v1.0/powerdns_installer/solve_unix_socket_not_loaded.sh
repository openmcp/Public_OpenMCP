#!/bin/bash
# https://askubuntu.com/questions/705458/ubuntu-15-10-mysql-error-1524-unix-socket
/etc/init.d/mysql stop
mysqld_safe --skip-grant-tables &
#mysql -uroot
#use mysql;

mysql -uroot -D mysql -e "update user set password=PASSWORD('ketilinux') where User='root'"
mysql -uroot -D mysql -e "update user set plugin='mysql_native_password'"

#quit;

/etc/init.d/mysql stop
kill -9 $(pgrep mysql)
/etc/init.d/mysql start
