include .env

build:
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o dnd_wild_magic_bot .

deploy: build
	ssh $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl stop dnd_wild_magic_bot'
	scp ./dnd_wild_magic_bot $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_PATH)/dnd_wild_magic_bot
	ssh $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl start dnd_wild_magic_bot'
