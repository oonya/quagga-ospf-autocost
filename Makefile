.PHONY: ocaca
ocaca:
	/usr/local/go/bin/go run cmd/ocaca/main.go

.PHONY: evalSeq
evalSeq:
	/usr/local/go/bin/go run cmd/evalution/main.go iperf3

.PHONY: evalRTT
evalRTT:
	/usr/local/go/bin/go run cmd/evalution/main.go ping
