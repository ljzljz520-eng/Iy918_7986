# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	guardpanel.local/guardpanel/cmd/guardpanel	0.002s
?   	guardpanel.local/guardpanel/internal/config	[no test files]
ok  	guardpanel.local/guardpanel/internal/domain	0.001s
ok  	guardpanel.local/guardpanel/internal/httpapi	0.008s
ok  	guardpanel.local/guardpanel/internal/importer	0.001s
ok  	guardpanel.local/guardpanel/internal/notifier	0.001s
ok  	guardpanel.local/guardpanel/internal/report	0.001s
--- FAIL: TestBusiness11Regression (0.01s)
    regression_test.go:38: expected explicit timeout, got sent after 1 attempts
FAIL
FAIL	guardpanel.local/guardpanel/internal/service	0.012s
ok  	guardpanel.local/guardpanel/internal/store	0.009s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/guardpanel): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/guardpanel): exit `0`
