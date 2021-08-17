.PHONY: start refresh watch stop


start:
	swag init >> /tmp/qnhd.log 2>&1 &
	air >> /tmp/qnhd.log 2>&1 &

refresh:
	swag init >> /tmp/qnhd.log 2>&1 &

watch:
	tail -f /tmp/qnhd.log

stop:
	pkill air
	rm /tmp/qnhd.log 2>/dev/null