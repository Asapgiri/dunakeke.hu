name=dunakeke.hu

build:
	go build -o ${name}

run:
	./${name}

reflex:
	reflex -R '\.git' -r '\.go' -s -- make build run

mongo-start:
	brew services start mongodb-community

mongo-stop:
	brew services stop mongodb-community

mongo-info:
	brew services info mongodb-community
