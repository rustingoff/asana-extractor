# ASANA Extractor

\***\*To run Golang app\*\***

`$ go run cmd/extractor/*.go --config <path_to_your_config_yaml_file>`

or

`$ go run cmd/extractor/*.go`
_without config argument and app will use default one (./config/config.yaml)_

### Configuration

_this is configuration structure, you can put your configuration data here_

```
app:
  api_auth_token: ""
  api_url: ""

  workspace_gid: ""

  user:
    path: "users"
    limit: 1 # 0 < limit <= 100

  project:
    path: "projects"
    limit: 1 # 0 < limit <= 100

```
