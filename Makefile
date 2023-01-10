.PHONY: init
init:
	sudo /usr/local/go/bin/go mod tidy

.PHONY: evalInit
evalInit:
	sudo /usr/bin/rm -f /home/oonya/quagga-ospf-autocost/rtt*
	sudo /usr/bin/rm -f /home/oonya/caputure/*

.PHONY: ocaca
ocaca:
	/usr/local/go/bin/go run cmd/ocaca/main.go

.PHONY: evalSeq
evalSeq:
	/usr/local/go/bin/go run cmd/evalution/sequence/main.go

.PHONY: evalRTT
evalRTT:
	/usr/local/go/bin/go run cmd/evalution/rtt/main.go

.PHONY: ocacaLog
ocacaLog:
	echo '' > /home/oonya/quagga-ospf-autocost/ocaca.log
	tail -f /home/oonya/quagga-ospf-autocost/ocaca.log

.PHONY: evalLog
evalLog:
	echo '' > /home/oonya/quagga-ospf-autocost/evalution.log
	tail -f /home/oonya/quagga-ospf-autocost/evalution.log
