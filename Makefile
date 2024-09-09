.PHONY: run_rolldice run_notification run_all

run_rolldice:
	go run ./cmd/rolldice/main.go

run_notification:
	go run ./cmd/notification/main.go
