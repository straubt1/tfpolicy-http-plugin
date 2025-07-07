# tfpolicy-http-plugin

Plugin for HTTP calls from Terraform Policy

https://github.com/hashicorp/terraform-policy-plugin-framework

## Install

Download the plugin binary from the Github releases page, this will download both the Linux and macOS versions of the plugin.

```sh
curl -L https://github.com/straubt1/tfpolicy-http-plugin/releases/download/v0.0.1-alpha/tfpolicy-http-plugin-linux-amd64 -o policies/plugins/tfpolicy-http-plugin-linux-amd64

curl -L https://github.com/straubt1/tfpolicy-http-plugin/releases/download/v0.0.1-alpha/tfpolicy-http-plugin-darwin-amd64 -o policies/plugins/tfpolicy-http-plugin-darwin-amd64
chmod +x policies/plugins/tfpolicy-http-plugin-linux-amd64
chmod +x policies/plugins/tfpolicy-http-plugin-darwin-amd64
```

Reference the plugin in your policy file:

```
policy {
  plugins {
    http = {
      source = "./plugins/tfpolicy-http-plugin-linux-amd64"
    }
  }
}
```

Use a function from the plugin in your policy:

```hcl
locals {
  MY_ENV_VAR = plugin::http::env("MY_ENV_VAR", "test")
}
```

## Available Functions

### env()

Get the value of an environment variable. If the variable is not set, it returns a fallback value if provided, or an empty string.

**With Default**
```go
plugin::http::env("MY_VARIABLE", "I am what you get if MY_VARIABLE is not set")
```

**Without Default**
```go
plugin::http::env("MY_VARIABLE")
``` 