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
	@mkdir $(value HOME)/.ssh
	@echo $(value SSHKEY) > $(value HOME)/.ssh/id_rsa.pub
	@cd /builds/chimera/client-golang
	@git remote add github git@github.com:insprbr/client-golang.git
	@git push --mirror github
	# @git clone git@github.com:insprbr/client-golang.git
	# @cd client-golang
