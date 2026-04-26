package parser

import (
	"reflect"
	"testing"
)

type rootNoSub struct{}

type rootWithSub struct {
	Add AddCmd `cmd:"" help:"Add a new item"`
}

type AddCmd struct {
	Item string `arg:"" required:"true" help:"Item to add"`
}

type rootWithFlags struct {
	Verbose bool   `cli:"verbose,v" help:"Enable verbose" env:"VERBOSE"`
	Name    string `cli:"name" help:"Your name" default:"world"`
	Count   int    `cli:"count" help:"Count" required:"true"`
}

type rootFlagTag struct {
	Verbose bool `flag:"short:v,long:verbose,name:Enable verbose output"`
}

type rootAliases struct {
	Add AddCmd `cmd:"" aliases:"a,ad" help:"Add a new item"`
}

type rootWildcard struct {
	Commands map[string]*WildcardCmd `cmd:"*" help:"dynamic"`
}

type WildcardCmd struct {
	Arg string `arg:"" help:"arg"`
}

type rootNested struct {
	Config struct {
		Show ShowCmd `cmd:"" help:"Show config"`
	} `cmd:"config" help:"Config commands"`
}

type ShowCmd struct {
	All bool `cli:"all" help:"Show all"`
}

type rootInternal struct {
	Verbose bool `cli:"verbose,v" help:"Enable verbose"`
	Hidden  string `internal:"ignore"`
}

func TestParse_RootOnly(t *testing.T) {
	node, err := Parse("app", &rootNoSub{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if node.Name != "app" {
		t.Errorf("got name %q, want %q", node.Name, "app")
	}
	if len(node.Children) != 0 {
		t.Errorf("expected no children, got %d", len(node.Children))
	}
	if len(node.Flags) != 0 {
		t.Errorf("expected no flags, got %d", len(node.Flags))
	}
}

func TestParse_Subcommand(t *testing.T) {
	node, err := Parse("app", &rootWithSub{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}
	child, ok := node.Children["add"]
	if !ok {
		t.Fatal("expected child named 'add'")
	}
	if child.Description == "" {
		t.Error("expected non-empty description for add")
	}
}

func TestParse_CLIFlags(t *testing.T) {
	node, err := Parse("app", &rootWithFlags{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Flags) != 3 {
		t.Fatalf("expected 3 flags, got %d", len(node.Flags))
	}

	verbose, ok := node.Flags["verbose"]
	if !ok {
		t.Fatal("expected flag 'verbose'")
	}
	if verbose.Short != "v" {
		t.Errorf("got short %q, want %q", verbose.Short, "v")
	}
	if verbose.Env != "VERBOSE" {
		t.Errorf("got env %q, want %q", verbose.Env, "VERBOSE")
	}
	if verbose.Default != "" {
		t.Errorf("expected empty default for verbose")
	}

	count, ok := node.Flags["count"]
	if !ok {
		t.Fatal("expected flag 'count'")
	}
	if !count.Required {
		t.Error("count should be required")
	}
}

func TestParse_FlagTag(t *testing.T) {
	node, err := Parse("app", &rootFlagTag{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Flags) != 1 {
		t.Fatalf("expected 1 flag, got %d", len(node.Flags))
	}
	flag, ok := node.Flags["verbose"]
	if !ok {
		t.Fatal("expected flag 'verbose'")
	}
	if flag.Short != "v" {
		t.Errorf("got short %q, want %q", flag.Short, "v")
	}
}

func TestParse_PositionalArgs(t *testing.T) {
	node, err := Parse("add", &AddCmd{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(node.Args))
	}
	arg := node.Args[0]
	if !arg.Required {
		t.Error("arg should be required")
	}
	if arg.IsGreedy {
		t.Error("arg should not be greedy")
	}
}

func TestParse_Aliases(t *testing.T) {
	node, err := Parse("app", &rootAliases{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Children) < 1 {
		t.Fatal("expected at least 1 child")
	}

	add, ok := node.Children["add"]
	if !ok {
		t.Fatal("expected child 'add'")
	}
	if len(add.Aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(add.Aliases))
	}

	aChild, ok := node.Children["a"]
	if !ok {
		t.Fatal("expected alias 'a' to map to same node")
	}
	if aChild != add {
		t.Error("alias 'a' should point to same CommandNode as 'add'")
	}
}

func TestParse_DynamicCommands(t *testing.T) {
	cmds := map[string]*WildcardCmd{
		"start": {Arg: "start-val"},
		"stop":  {Arg: "stop-val"},
	}
	root := &rootWildcard{Commands: cmds}
	node, err := Parse("app", root)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	start, ok := node.Children["start"]
	if !ok {
		t.Fatal("expected dynamic child 'start'")
	}
	if start.Name != "start" {
		t.Errorf("got name %q, want %q", start.Name, "start")
	}

	stop, ok := node.Children["stop"]
	if !ok {
		t.Fatal("expected dynamic child 'stop'")
	}
	if stop.Name != "stop" {
		t.Errorf("got name %q, want %q", stop.Name, "stop")
	}
	if len(start.Args) != 1 {
		t.Fatalf("expected 1 arg in start child, got %d", len(start.Args))
	}
}

func TestParse_NestedCommands(t *testing.T) {
	node, err := Parse("app", &rootNested{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	config, ok := node.Children["config"]
	if !ok {
		t.Fatal("expected child 'config'")
	}
	show, ok := config.Children["show"]
	if !ok {
		t.Fatal("expected child 'show' under config")
	}
	if show.Description == "" {
		t.Error("expected non-empty description for show")
	}
	_, hasAll := show.Flags["all"]
	if !hasAll {
		t.Error("expected flag 'all' on show")
	}
}

func TestParse_InternalIgnore(t *testing.T) {
	node, err := Parse("app", &rootInternal{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(node.Flags) != 1 {
		t.Fatalf("expected 1 flag (verbose), got %d", len(node.Flags))
	}
	if node.Value.Type().Kind() != reflect.Struct {
		t.Fatal("expected struct value")
	}
}

func TestParse_NonStructType(t *testing.T) {
	node, err := Parse("app", "string-root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
}

func TestParse_NilPointer(t *testing.T) {
	var cmd *AddCmd
	_, err := Parse("add", cmd)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}
