.PHONY: start refresh watch stop tstart

tstart:
	go run main.go

tstop:
	kill -9 

start:
	air >> /tmp/qnhd.log 2>&1 &

watch:
	tail -f /tmp/qnhd.log

stop:
	pkill air
	rm /tmp/qnhd.log 2>/dev/null