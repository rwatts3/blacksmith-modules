#
# build is a shortcut to build the Go modules. Since there is no main packages,
# this step is only here to tidy the dependencies across modules.
#
build:
	goreleaser --snapshot --rm-dist

#
# test is a shortcut to run the Go tests across modules.
#
test:
	sh ./scripts/test.sh
