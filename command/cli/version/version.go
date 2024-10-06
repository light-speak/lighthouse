// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package version

import (
	"fmt"

	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/version"
)


type Version struct {}

func (c *Version) Name() string {
	// Func:Name user code start. Do not remove this comment.
	return "app:version"
	// Func:Name user code end. Do not remove this comment. 
}

func (c *Version) Usage() string {
	// Func:Usage user code start. Do not remove this comment.
	return "Show the version of lighthouse"
	// Func:Usage user code end. Do not remove this comment. 
}

func (c *Version) Args() []*command.CommandArg {
	return []*command.CommandArg{
		// Func:Args user code start. Do not remove this comment.
		
		// Func:Args user code end. Do not remove this comment. 
	}
}

func (c *Version) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		// Func:Action user code start. Do not remove this comment.
		printVersion()
		// Func:Action user code end. Do not remove this comment. 
		return nil
	}
}


// Section: user code section start. Do not remove this comment.


func printVersion() {
	fmt.Printf("Version: %s\n", version.Version)
}


// Section: user code section end. Do not remove this comment. 