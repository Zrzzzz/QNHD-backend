start:
	swag init >> /tmp/qnhd.log 2>&1 & && air >> /tmp/qnhd.log 2>&1 &

refresh:
	swag init >> /tmp/qnhd.log 2>&1 &

stop:
	ps -ef | grep air | awk {'print $2'} | xargs kill -9
	rm /tmp/qnhd.log 2>/dev/null