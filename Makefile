LINTER_VERSION = v2.1.5

lint:
	@if ! bin/golangci-lint --version | grep -q $(LINTER_VERSION); \
		then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(LINTER_VERSION); fi;
	bin/golangci-lint run --fix