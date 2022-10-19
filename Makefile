test:
	go test ./... -v -count=1

CURRENT_TAG=$(shell git tag --sort version:refname | tail -n 1)

first_release:
	git tag v0.0.1
	git push origin v0.0.1

rev_release:
	git tag $(shell exoskeleton rev -i $(CURRENT_TAG))
	git push origin $(shell exoskeleton rev -i $(CURRENT_TAG))

release:
	@if [[ "$(CURRENT_TAG)" == "" ]]; then $(MAKE) first_release; else $(MAKE) rev_release; fi