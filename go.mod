module github.com/tokenized/teller_client

go 1.14

require (
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/tokenized/channels v0.0.0-20220902192544-9fd058a84b32
	github.com/tokenized/envelope v1.0.1-0.20220902162954-c431b788f500
	github.com/tokenized/logger v0.1.0
	github.com/tokenized/pkg v0.4.1-0.20221111193511-06a224463802
	github.com/tokenized/specification v1.1.2-0.20220902163651-058d13f0a70f
	github.com/tokenized/threads v0.1.1-0.20220902155404-d844f8ac41b5
)

replace github.com/tokenized/specification => ../specification
