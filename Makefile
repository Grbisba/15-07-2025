CONFIG_FILE_PATH="./infra/config.json"
.PHONY: lines
lines:
	git ls-files | xargs wc -l

.PHONY: build
build:
	go build -o app.o ./cmd/app/

.PHONY: run
run: build
	env CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) ./app.o

.PHONY: key-gen
key-gen: secret

secret: create-dir gen-private-key gen-public-key

create-dir:
	@mkdir -p certs
gen-private-key:
	@openssl genrsa -out certs/private_key.pem 2048
gen-public-key:
	@openssl rsa -in certs/private_key.pem -pubout -outform PEM -out certs/public_key.pem
