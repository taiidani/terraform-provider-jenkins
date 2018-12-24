build:
	docker run --rm -it \
		--volume "$(shell pwd):/srv/jekyll" \
		--volume "$(shell pwd)/.bundles:/usr/local/bundle" \
		--env PAGES_REPO_NWO="taiidani/terraform-provider-jenkins" \
		-p "4000:4000" \
		jekyll/jekyll jekyll serve --watch
