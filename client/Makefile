EXEC_FILE = client
POLICY_EXEC_FILE = ./bin/$(EXEC_FILE)
CLIENT = /mnt/d/workspace/falplus/common:/mnt/d/workspace/falplus/dependency:/mnt/d/workspace/falplus/client

.PHONY: all clean home

all: clean $(POLICY_EXEC_FILE)

$(POLICY_EXEC_FILE):
	@echo "all ..."
	find ./src -name "*.go"  |xargs gofmt -w
	go build -v $(EXEC_FILE)
	@mkdir -p ./bin
	@mv $(EXEC_FILE) ./bin/

clean:
	@echo "clean ..."
	rm -f $(POLICY_EXEC_FILE)

install:
	@echo "install ..."
	install $(POLICY_EXEC_FILE) /usr/local/bin/$(EXEC_FILE)

home: export GOPATH=$(CLIENT)
home:
	@echo "all ..."
	find ./src -name "*.go"  |xargs gofmt -w
	go build -v $(EXEC_FILE)
	@mkdir -p ./bin
	@mv $(EXEC_FILE) ./bin/