# Automation of building the GoMenacing protobuf schema.
#
PROTOC_OUTDIR ?= ./pkg/gomschema
PROTOC_INCDIR ?= ./api/gomschema
PROTOC_ARGS   ?= -I "$(PROTOC_INCDIR)"
PROTOC_LANGS  ?= --go_out=$(PROTOC_OUTDIR) --python_out=$(PROTOC_OUTDIR) --java_out=$(PROTOC_OUTDIR) --cpp_out=$(PROTOC_OUTDIR)
PROTOC_SCHEMA ?= ./api/gomschema/gomschema.proto

GOPATH        ?= ${HOME}/go

PROTOC_CMD    ?= protoc

DOCKER_IMAGE ?= kfsone/gomprotoc
IMAGE_VER    ?= latest

all:
	@echo "Compile gomschema protobuffers schema."
	@echo ""
	@echo "Use:"
	@echo "  make clean    -- remove the existing output."
	@echo "  make protoc   -- use the (installed) protoc compiler directly."
	@echo "  make wsl      -- use Windows Subsystem for Linux (win only)."
	@echo "  make inwsl    -- you're inside WSL already (win only)."
	@echo "  make deb      -- run from a debian install/wsl."
	@echo "  make docker   -- use a Docker container to compile with."
	@echo ""
	@echo "Use PROTOC_LANGS to override the default languages ($PROTOC_LANGS)"

.PHONY: wsl
wsl:
	wsl -- bash -c "$(MAKE) deb"

.PHONY: deb
deb:
	sudo apt update && sudo apt install --upgrade git golang protobuf-compiler && \
		mkdir -p "${GOPATH}/bin" && \
		export PATH="${PATH}:${GOPATH}/bin" && \
		go get -u google.golang.org/protobuf/cmd/protoc-gen-go && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go && \
		$(MAKE) protoc

docker-image: Dockerfile
	docker build --tag "$(DOCKER_IMAGE):$(IMAGE_VER)" .

docker-publish: docker-image
	docker push "$(DOCKER_IMAGE):$(IMAGE_VERSION)"

.PHONY: docker
docker:
	docker run --rm -v $(PWD):/gom $(DOCKER_IMAGE)

.PHONY: protoc
protoc:
	"$(PROTOC_CMD)" \
			$(PROTOC_ARGS) \
			$(PROTOC_LANGS) \
			$(PROTOC_SCHEMA) && \
		echo "Done."

