unexport GOBIN

NAME         	= pet-adoption
GO           	= go
SPECE_DIR       = ./_rpmbuild
RPMBUILD_DIR	= /root/rpmbuild
VERSION = 1.0.0
RELEASE = 15
DOCKERREPO = 192.168.1.100:5000/gzsunrun
BRANCH ?= master

all: build

uppkg:
	@echo ">> update pkg..."
	@glide --home /root/.glide/$(NAME)/ up

pkg:
	@echo ">> get pkg..."
	glide --home /root/.glide/$(NAME)/  install

build:
	@echo ">> building code..."
	go build -ldflags "-s -w -X main._VERSION_=$(VERSION)-$(RELEASE)" -o bin/$(NAME)

rpmbuild:
	@echo ">> rpm building"
	mkdir -p $(RPMBUILD_DIR)/BUILD/$(NAME)/usr/local/bin/
	mkdir -p $(RPMBUILD_DIR)/BUILD/$(NAME)/etc/systemd/system/
	mkdir -p $(RPMBUILD_DIR)/BUILD/$(NAME)/etc/pet-adoption/

	cp ./bin/$(NAME) -f $(RPMBUILD_DIR)/BUILD/$(NAME)/usr/local/bin/
	cp $(SPECE_DIR)/$(NAME).service -f $(RPMBUILD_DIR)/BUILD/$(NAME)/etc/systemd/system/
	cp $(SPECE_DIR)/config.yml -f $(RPMBUILD_DIR)/BUILD/$(NAME)/etc/pet-adoption/
	cp $(SPECE_DIR)/$(NAME).spec -f $(RPMBUILD_DIR)/SPECS/
	sed -i 's/^Version:.*/Version:$(VERSION)/' $(RPMBUILD_DIR)/SPECS/$(NAME).spec
	sed -i 's/^Release:.*/Release:$(RELEASE).el7/' $(RPMBUILD_DIR)/SPECS/$(NAME).spec
	rpmbuild -bb $(RPMBUILD_DIR)/SPECS/$(NAME).spec

dockerbuild:
	@echo ">> build docker image"
	docker build -t $(DOCKERREPO)/$(NAME):latest .
	docker tag $(DOCKERREPO)/$(NAME):latest $(DOCKERREPO)/$(NAME):$(VERSION)-$(RELEASE)

pimage:
	@echo ">> push release docker image"
	docker push $(DOCKERREPO)/$(NAME):$(VERSION)-$(RELEASE)
	docker push $(DOCKERREPO)/$(NAME):latest

ciimage:
	@echo ">> build dev docker image"
	docker build -t $(DOCKERREPO)/$(NAME):$(BRANCH) .
	@echo ">> push dev docker image"
	docker push  $(DOCKERREPO)/$(NAME):$(BRANCH)