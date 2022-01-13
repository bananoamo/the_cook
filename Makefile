

NAME=./main.go

all:
	go build .
	./the_cook

run:
	go run .

clean:
	rm -f main

.PHONEY: all run clean