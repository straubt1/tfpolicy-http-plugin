# tfpolicy-http-plugin

Plugin for HTTP calls from Terraform Policy

https://github.com/hashicorp/terraform-policy-plugin-framework

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