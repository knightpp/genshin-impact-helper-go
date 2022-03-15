
default:
	just --list

publish:
	export BUILDAH_FORMAT="docker" 
	heroku container:push web
	heroku container:release web