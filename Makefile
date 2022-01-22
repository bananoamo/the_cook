.PHONEY: all build clean

NAME=the_cook

all: build
	./$(NAME)

build:
	go build .

clean:
	rm -f $(NAME)
