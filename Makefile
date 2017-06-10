NAME=go-chapter-route-recommendation
SOURCE=main.go

all: clean depend build

pre_commit_hook:
	cp $(CURDIR)/tools/git-pre-commit-hook.sh $(CURDIR)/.git/hooks/pre-commit
	chmod +x $(CURDIR)/.git/hooks/pre-commit

clean:
	rm -rf build/

depend:
	go get -u -v github.com/Masterminds/glide
	glide install

build:
	go build -o build/$(NAME) $(SOURCE)

fmt:
	go fmt $(shell glide novendor)

vet:
	go vet $(shell glide novendor)

lint:
	go get -u -v github.com/golang/lint/golint
	for file in $(shell find . -name '*.go' -not -path './vendor/*'); do golint $${file}; done
