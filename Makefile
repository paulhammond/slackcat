prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/slackcat; then rm -rf src/github.com/whosonfirst/slackcat; fi
	mkdir -p src/github.com/whosonfirst/slackcat
	cp slackcat.go src/github.com/whosonfirst/slackcat/

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin:	self fmt
	go build -o bin/slackcat cmd/slackcat.go