CC=go

gotodo: gotodo.go
	$(CC) build ./gotodo.go

clean:
	rm -rf gotodo
