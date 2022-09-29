deps:
	pip install pylint
	wget https://github.com/golangci/golangci-lint/releases/download/v1.49.0/golangci-lint-1.49.0-linux-amd64.tar.gz -O golangci-lint.tar.gz
	tar xvf golangci-lint.tar.gz
	mv golangci-lint-1.49.0-linux-amd64/golangci-lint ${GOPATH}/bin/golangci-lint-1.49
	rm -rf golangci*

clean:
	rm -rf build dist rbfopt_go.egg-info
	rm -rf /tmp/rbf*

test:
	go test -tags testing -count=1 -v ./...

lint:
	golangci-lint-1.49 run ./...
	pylint ./rbfoptgo

python_release_build:
	rm -rf ./build/* ./dist/*
	python setup.py sdist bdist_wheel
	twine check dist/*

python_release_test:
	twine upload --repository-url https://test.pypi.org/legacy/ dist/*

python_release_prod:
	twine upload dist/*

deps_fedora:
	sudo dnf install -y coin-or-Bonmin