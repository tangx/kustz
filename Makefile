
tidy:
	go mod tidy

build:
	go build ./cmd/kustz

test.deployment:
	make test TARGET=Test_KustzDeployment
test.kustomize:
	make test TARGET=Test_KustzKustomize


test:
	go test -timeout 30s -run ^$(TARGET) github.com/tangx/kustz/pkg/kustz -v -count=1
	