build:
	go build -o server ./server
	go build -o cli ./cli

release:
	git tag v1.0.0
	git push origin v1.0.0