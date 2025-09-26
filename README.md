# mcp262

Model Context Protocol (MCP) server exposing Test262 test metadata, code, harness files, ECMAScript spec search and test runner integration for the YavaShark JavaScript engine (and other engines via pluggable providers).

## Features
- Test inventory: count, list (recursive / non‑recursive), pagination.
- Test status & output retrieval (from last CI / runner results).
- Filter by status (PASS, FAIL, SKIP, TIMEOUT, CRASH, PARSE_ERROR, NOT_IMPLEMENTED, RUNNER_ERROR).
- Search: repository paths, tests, spec text, spec section titles.
- Code access & in‑memory editing: individual test files, harness files, harness set per test.
- Spec utilities: fetch spec section, intrinsic lookup (%Array.prototype.toString%).
- Parallel runner integration (configurable workers) prepared for large scale test execution.
- Modular provider interfaces for: Test data, Code (tests + harness), Spec content.
- HTTP Streamable MCP server (default :8080) compatible with MCP clients.

## High Level Architecture
```
+------------------+        +------------------+
|  MCP Client(s)   | <----> |  mcp262 Server   |
+------------------+        +---------+--------+
                                      |
          +---------------------------+---------------------------+
          |                           |                           |
   TestProvider                 TestCodeProvider             SpecProvider
 (statuses, lists,             (test & harness code)       (ECMA spec text,
  search, outputs)                                          search, intrinsic)

                Runner (parallel execution, results aggregation)
```
Providers are set at startup (see main in mcp262.go) and can be swapped to target other engines or data sources.

## Getting Started
### Prerequisites
- Go 1.24+
- Test262 repository (or a directory mirroring its structure) available locally.

### Install
```
git clone https://github.com/Sharktheone/mcp262.git
cd mcp262
go build ./...
```

### Configuration
You can configure via:
1. TOML file (default: config.toml) – see config.example.toml
2. Environment variables
3. CLI flags (override both file & env)

Config keys:
- repo_path (REPO_PATH / --repo) : path to external repository root (default ./)
- test_root_dir (TEST_ROOT_DIR / --test_root) : root to test262 tests (default ./test262/test)
- workers (WORKERS / --workers) : parallel workers for runner (default 256)

Example config.toml:
```
workers = 128
repo_path = "./"
test_root_dir = "./test262/test"
```
Run with explicit config file:
```
go run . --config config.toml
```
Or override:
```
WORKERS=64 TEST_ROOT_DIR=./test262/test go run . --workers 32
```

### Running the Server
```
go run .
```
Server listens (by default) on 0.0.0.0:8080 providing MCP over a streamable HTTP endpoint.

## MCP Tools Overview
(Category / Tool Name -> Purpose)
- Tests
  - NumTestsTotal, NumTestsInDir, NumTestsInDirRecursive
  - GetTestsInDir, GetTestsInDirRecursive
  - GetTestStatus, GetTestStatusesInDir, GetTestStatusesInDirRecursive
  - GetTestsWithStatusInDir, GetTestsWithStatusInDirRecursive
  - GetFailedTestsInDir, GetFailedTestsInDirRecursive
  - GetTestOutput
  - SearchDir, SearchDirIn, SearchTest, SearchTestInDir
- Code / Harness
  - GetTestCode, GetHarnessForTest, GetHarness, GetHarnessCode
  - GetHaressFiles (sic), GetHarnessFilesForTest
  - SetTestCode, SetHarnessCode, ResetEdits
- Spec
  - GetSpec, SpecForIntrinsic, SearchSpec, SearchSections
- Runner (if added in tools.AddRunnerTools) – execution & diff related utilities (see runner directory).

Pagination fields: page, page_size, returned, remaining, total.

## Runner
The runner package (runner/) manages parallel execution of tests (Workers) and stores summarized results accessible to TestProvider implementations. Configure concurrency via workers.


## License
MIT – see LICENSE.

## Disclaimer
mcp262 is experimental; APIs and tool names may change.
