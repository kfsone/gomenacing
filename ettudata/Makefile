# This is an awful but obvious place to record how to generate the pb.go file.
#
FLATC_OUTDIR ?= .
FLATC_INCDIR ?= flatbufs
FLATC_ARGS   ?= --natural-utf8 --scoped-enums --size-prefixed
FLATC_GENS   ?= --gen-mutable --gen-nullable --gen-generated
FLATC_LANGS  ?= --go --python --csharp --rust --js --ts
FLATC_SCHEMA ?= flatbufs/ettuschema.fbs

FLATC_CMD    ?= flatc

DOCKER_IMAGE ?= kfsone/ettuflatc
IMAGE_VER    ?= latest

all:
	@echo "Compile ettudata flatbuffers schema."
	@echo ""
	@echo "Use:"
	@echo "  make clean    -- remove the existing output."
	@echo "  make flatc    -- use the flatc compiler directly."
	@echo "  make wsl      -- use Windows Subsystem for Linux (win only)."
	@echo "  make docker   -- use a Docker container to compile with."
	@echo ""
	@echo "Use FLATC_LANGS to override the default languages ($FLATC_LANGS)"

.PHONY: wsl
wsl:
	wsl -- bash -c "make inwsl"

.PHONY: inwsl
inwsl:
	sudo apt update && apt install --upgrade flatbuffers-compiler
	$(MAKE) flatc

docker-image: Dockerfile
	docker build --tag "$(DOCKER_IMAGE):$(IMAGE_VER)" .

docker-publish: docker-image
	docker push "$(DOCKER_IMAGE):$(IMAGE_VERSION)"

.PHONY: docker
docker:
	docker run --rm -v $(PWD):/ettudata kfsone/ettuflatc

.PHONY: flatc
flatc:
	@"$(FLATC_CMD)" \
			-o "$(FLATC_OUTDIR)" \
			-I "$(FLATC_INCDIR)" \
			$(FLATC_ARGS) \
			$(FLATC_GENS) \
			$(FLATC_LANGS) \
			$(FLATC_SCHEMA) && \
		echo "Done."

