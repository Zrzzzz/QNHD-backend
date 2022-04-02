.PHONY: run start stop

run:
	go run .

start:
	go build -o backend .
	nohup ./backend &

stop:
	pkill backend
	rm nohup.out