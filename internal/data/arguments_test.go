package data

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

type parseArgumentsTest struct {
	when     string
	args     []string
	expected *Arguments
}

var parseArgumentsTests = []parseArgumentsTest{
	{
		when: "provider and region is defined",
		args: []string{"autovpn", "linode", "eu-central-1"},
		expected: &Arguments{
			Command:  Default,
			Provider: "linode",
			Region:   "eu-central-1",
		},
	},
	{
		when: "listing providers",
		args: []string{"autovpn", "list"},
		expected: &Arguments{
			Command: ListArgs,
		},
	},
	{
		when: "listing regions",
		args: []string{"autovpn", "linode", "list"},
		expected: &Arguments{
			Command: ListArgs,
			Provider: "linode",
		},
	},
	{
		when: "listing zombies",
		args: []string{"autovpn", "linode", "zombies"},
		expected: &Arguments{
			Command:  ListZombies,
			Provider: "linode",
		},
	},
	{
		when: "purging",
		args: []string{"autovpn", "linode", "purge"},
		expected: &Arguments{
			Command:  Purge,
			Provider: "linode",
		},
	},
	{
		when: "listing all zombies",
		args: []string{"autovpn", "zombies"},
		expected: &Arguments{
			Command: ListZombies,
		},
	},
	{
		when: "purging all",
		args: []string{"autovpn", "purge"},
		expected: &Arguments{
			Command: Purge,
		},
	},
	{
		when: "only provider defined",
		args: []string{"autovpn", "linode"},
		expected: &Arguments{
			Command: Usage,
		},
	},
	{
		when: "no argument is defined",
		args: []string{"autovpn"},
		expected: &Arguments{
			Command: Usage,
		},
	},
	{
		when: "help flag is defined",
		args: []string{"autovpn", "--help"},
		expected: &Arguments{
			Command: Usage,
		},
	},
	{
		when: "version flag is defined",
		args: []string{"autovpn", "--version"},
		expected: &Arguments{
			Command: Usage,
		},
	},
}

func TestParseArguments(t *testing.T) {
	for _, test := range parseArgumentsTests {
		os.Args = test.args
		actual, _ := ParseArguments()
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("When %s: actual %v not equal to expected %v", test.when, actual, test.expected)
		}
		// avoid errors generated by flag when not resetting parsed flags
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}
}