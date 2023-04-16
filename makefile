CONFIG_PATH=${HOME}/.logservice/

.PHONY: gproto
gproto:
	protoc api/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--proto_path=.

.PHONY: test
test:
	go test -race ./...

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem \
	-ca-key=ca-key.pem \
	-config=test/ca-config.json \
	-profile=server \
	test/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca.pem \
	-ca-key=ca-key.pem \
	-config=test/ca-config.json \
	-profile=client \
	test/client-csr.json | cfssljson -bare client
	mv *.pem *csr ${CONFIG_PATH}

.PHONY: invalidgencert
invalidgencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert -initca test/ca-csr-invalid.json | cfssljson -bare ca-invalid
	cfssl gencert -ca=ca.pem \
	-ca-key=ca-key.pem \
	-config=test/ca-config.json \
	-profile=server \
	test/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca-invalid.pem \
	-ca-key=ca-invalid-key.pem \
	-config=test/ca-config-invalid.json \
	-profile=client \
	test/client-csr-invalid.json | cfssljson -bare client
	mv *.pem *csr ${CONFIG_PATH}