GOCMD=go
GOINSTALL=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run

install: 
	$(GOINSTALL) ${GOPATH}/src/github.com/suhay/quest-log

clean: 
	$(GOCLEAN)
	rm -f ${GOPATH}/bin/quest-log

gqlgen:
	$(GORUN) github.com/99designs/gqlgen -v

dev:
	CompileDaemon -directory=. -color=true -command="./quest-log"

test:
	LOCAL_PATH=fixtures go test -v github.com/suhay/quest-log/tests