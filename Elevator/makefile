BUILDPATH = $(CURDIR)
GO = $(shell which go)
GOINSTALL=$(GO) install
GOCLEAN = $(GO) clean
GOBUILD = $(GO) build

export GOPATH = ${CURDIR}
export GOBIN=${CURDIR}/bin

myname:
	@echo "Call to makefile"

install:
	@$(GOINSTALL) main

execute:
	@echo "Execute"
	@echo "---------------"
	@./bin/main

makedir:
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	@if [ ! -d $(BUILDPATH)/pkg ] ; then mkdir -p $(BUILDPATH)/pkg ; fi

clean: 
	@echo "Cleaned: /bin and /pkg deleted"
	@rm -rf $(BUILDPATH)/bin
	@rm -rf $(BUILDPATH)/pkg

run: myname makedir install execute



