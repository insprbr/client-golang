TAG :="$(shell git describe)-ci"

IMG_REPO:=gcr.io/red-inspr/clients/

MODULE_NAME:=client-golang

.PHONY: build
build:
	@go build -o $(MODULE_NAME) -tags musl

.PHONY: lint
lint:
	@golint -set_exit_status ./...

.PHONY: deploy
deploy:
	@eval "$(ssh-agent -s)"
	@echo "$(SSHKEY)"
	@git clone git@github.com:insprbr/client-golang.git
	@cd client-golang
