package directives

import "fmt"

func ExampleParse() {
	input := "::blog-draft Write about testing"
	result := Parse(input, ParserConfig{Source: "session-abc"})

	fmt.Printf("%s|%s|%s\n", result.Directives[0].Command, result.Directives[0].Prompt, result.Directives[0].Source)
	// Output: blog-draft|Write about testing|session-abc
}

func ExampleParse_configCascade() {
	input := `::config voice=technical, project=carrier
::context_start Sprint
::blog-draft Post about sprint
::zoom Cron details
::config voice=casual
::note Funny bug
::zoom out
::adr Cron decision
::context_end`

	result := Parse(input, ParserConfig{})

	fmt.Printf("%s:%s\n", result.Directives[0].Command, result.Directives[0].Config["voice"])
	fmt.Printf("%s:%s\n", result.Directives[1].Command, result.Directives[1].Config["voice"])
	fmt.Printf("%s:%s\n", result.Directives[2].Command, result.Directives[2].Config["voice"])
	// Output:
	// blog-draft:technical
	// note:casual
	// adr:technical
}
