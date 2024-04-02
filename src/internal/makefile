MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
SRCS := $(shell find . -type f -name '*.go' ! -name '*test.go' ! -name '*mocks')
mocks: $(SRCS)
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done