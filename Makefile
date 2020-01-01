GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run

install: 
	$(GOINSTALL) $(shell pwd)/quest-logp.go

clean: 
	$(GOCLEAN)
	rm -f ${GOPATH}/bin/quest-log

gqlgen:
	$(GORUN) github.com/99designs/gqlgen -v

dev:
	CompileDaemon -directory=. -color=true -command="./quest-log"