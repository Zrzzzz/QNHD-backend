.PHONY: start refresh watch stop tstart

tstart:
	export GIN_MODE="debug"
	go run main.go

tstop:
	kill -9 

start:
	export GIN_MODE="release"
	air >> /tmp/qnhd.log 2>&1 &

watch:
	tail -f /tmp/qnhd.log

stop:
	pkill air
	rm /tmp/qnhd.log 2>/dev/null
