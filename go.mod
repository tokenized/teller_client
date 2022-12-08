module github.com/tokenized/teller_client

go 1.14

require (
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/tokenized/channels v0.0.0-20220902192544-9fd058a84b32
	github.com/tokenized/envelope v1.0.1-0.20220902162954-c431b788f500
	github.com/tokenized/logger v0.1.2-0.20221123201255-3047489d4997
	github.com/tokenized/pkg v0.4.1-0.20221202154605-c39f9adc94f5
	github.com/tokenized/specification v1.1.2-0.20221208151154-bb90ecbdf399
	github.com/tokenized/threads v0.1.1-0.20220902155404-d844f8ac41b5
)

replace github.com/tokenized/specification => ../specification
