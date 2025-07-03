// main.go

package main

import "github.com/hashicorp/terraform-policy-plugin-framework/policy-plugin/plugins"

func main() {
	plugins.RegisterFunction("echo", func(input string) (string, error) {
		// simple echo function, just return the string input directly.
		return input, nil
	})
	plugins.Serve()
}
