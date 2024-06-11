default: build

OAPI := ./v1/generated/model/model.yaml
OAPI += ./v1/generated/common/routes.yaml
OAPI += ./v1/generated/service/routes.yaml

build:
	@for dir in $(OAPI); do \
		echo "$$dir"; \
		oapi-codegen -config $$dir oapi.yaml; \
	done
