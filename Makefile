# Bootstrap: copy bootstrap-files to project root and replace app name / owner name placeholders.
# Usage: make bootstrap OWNER=<owner> APP=<app>
# Example: make bootstrap OWNER=myorg APP=mycli

.PHONY: bootstrap

bootstrap:
	@if [ -z "$(OWNER)" ] || [ -z "$(APP)" ]; then \
		echo "Usage: make bootstrap OWNER=<ownername> APP=<appname>"; \
		echo "Example: make bootstrap OWNER=myorg APP=mycli"; \
		exit 1; \
	fi
	@APP_UPPER=$$(echo $(APP) | tr '[:lower:]' '[:upper:]'); \
	mkdir -p ./.bootstrap; \
	cp Makefile ./.bootstrap/Makefile; \
	cp README.md ./.bootstrap/README.md; \
	cp -r bootstrap-files/. .; \
	mv bootstrap-files/ ./.bootstrap/bootstrap-files/; \
	mv cmd/appname cmd/$(APP); \
	mv data/logrotate.d/appname data/logrotate.d/$(APP); \
	echo "Replacing placeholders..."; \
	find . -type f ! -path './.bootstrap/*' -exec sed -i "s/APPNAME/$$APP_UPPER/g" {} \;; \
	find . -type f ! -path './.bootstrap/*' -exec sed -i "s/appname/$(APP)/g" {} \;; \
	find . -type f ! -path './.bootstrap/*' -exec sed -i "s/ownername/$(OWNER)/g" {} \;; \
	echo "Initializing git repository..."; \
	git init; \
	echo "Getting required dependencies..."; \
	go get -u github.com/urfave/cli/v2; \
	go get -u github.com/BurntSushi/toml; \
	go get -u github.com/joho/godotenv; \
	echo "Bootstrap complete."; \
	rm -rf ./Makefile; \
	mv Makefile.project Makefile; \
	rm -f ./.git/index; \
	git reset; \
	echo "Please remove the '.bootstrap' directory before proceeding."; \
	echo "If you need to run the bootstrap again, cd into the '.bootstrap' directory then run 'make clean' to reset the project state first"; \
	echo "--------------------------------"; \
	echo "--------------------------------"; \

clean:
	@if [ "$$(basename $$(pwd))" != ".bootstrap" ]; then \
		echo "clean can only be run from a directory named '.bootstrap'"; \
		exit 1; \
	fi; \
	rm -rf ../cmd/; \
	rm -rf ../internal/; \
	rm -rf ../data/; \
	rm -rf ../dist/; \
	rm -rf ../builds/; \
	rm -rf ../docs/; \
	rm -rf ../Makefile.project; \
	rm -rf ../go.mod; \
	rm -rf ../go.sum; \
	rm -rf ../.gitignore; \
	mv ./Makefile ../Makefile; \
	mv ./README.md ../README.md; \
	rm -rf ../bootstrap-files; \
	mv ./bootstrap-files ../bootstrap-files; \
	rm -f ../.git/index; \
	git reset; \
	rm -rf ../.bootstrap/; \
	echo "Clean complete."; \
