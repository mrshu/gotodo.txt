CC=go
INSTALL=install -c
DEST=/usr/bin

gotodo: gotodo.go
	$(CC) build ./gotodo.go

install: gotodo
	$(INSTALL) gotodo $(DEST)/gotodo

clean:
	rm -rf gotodo
