BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOBUILD=$(GO) build

EXECNAME=testexe


test :
        @echo $(BUILDPATH)

build :
        @echo "Build starting......."
        @if [ ! -d $(BUILDPATH)/bin ] ; then mkdir $(BUILDPATH)/bin ; fi
        $(GOBUILD) -o $(BUILDPATH)/bin/$(EXECNAME)
        @echo "Build successful......."
clean :
        @rm -rf  $(BUILDPATH)/bin
