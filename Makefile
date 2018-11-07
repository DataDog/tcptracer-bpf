# TODO: Move this to a Rakefile, for consistency with the other agent
DEBUG=1
UID=$(shell id -u)
PWD=$(shell pwd)

DOCKER_FILE?=Dockerfile
DOCKER_IMAGE?=datadog/tcptracer-bpf-builder

# If you can use docker without being root, you can do "make SUDO="
SUDO=$(shell docker info >/dev/null 2>&1 || echo "sudo -E")

# Generate and install eBPF program via gobindata
all: build-docker-image build-ebpf-object install-generated-go

build-docker-image:
	$(SUDO) docker build -t $(DOCKER_IMAGE) -f $(DOCKER_FILE) .

build-ebpf-object: build-docker-image
	$(SUDO) docker run --rm -e DEBUG=$(DEBUG) \
		-e CIRCLE_BUILD_URL=$(CIRCLE_BUILD_URL) \
		-v $(PWD):/src:ro \
		-v $(PWD)/ebpf:/dist/ \
		--workdir=/src \
		$(DOCKER_IMAGE) \
		make -f ebpf.mk build
	sudo chown -R $(UID):$(UID) ebpf

install-generated-go:
	cp ebpf/tcptracer-ebpf.go pkg/tracer/tcptracer-ebpf.go

delete-docker-image:
	$(SUDO) docker rmi -f $(DOCKER_IMAGE)

lint:
	./tools/lint -ignorespelling "agre " -ignorespelling "AGRE " .
	./tools/shell-lint .

# Build & run dockerized `nettop` command for testing 
# $ make all run-nettop
run-nettop:
	sudo docker build -t "tcptracer-bpf-dd-nettop" . -f tests/Dockerfile-nettop
	sudo docker run \
		--net=host \
		--cap-add=SYS_ADMIN \
		--privileged \
		-v /sys/kernel/debug:/sys/kernel/debug \
		tcptracer-bpf-dd-nettop

test:
	go list ./... | grep -v vendor | sudo -E PATH=${PATH} GOCACHE=off xargs go test -tags 'linux_bpf'

# TODO: Add linux_bpf tag so it runs CI tests w/ eBPF enabled
ci-test: build-ebpf-object
	go list ./... | grep -v vendor | sudo -E PATH=${PATH} GOCACHE=off xargs go test -tags ''
