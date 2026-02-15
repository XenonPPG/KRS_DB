update-deps:
	cd app && set GOPROXY=direct && set GOSUMDB=off && go clean -modcache
	cd app && set GOPROXY=direct && set GOSUMDB=off && go get -u github.com/XenonPPG/KRS_CONTRACTS@master
	cd app && go mod tidy

upload-image:
	docker build -t xenonppg/krs_db:latest ./app
	docker push xenonppg/krs_db:latest

.PHONY: update-deps, upload-image