build:
	go build -o bin/server ./server
	go build -o bin/cli ./cli

release:
	git tag v1.0.0
	git push origin v1.0.0