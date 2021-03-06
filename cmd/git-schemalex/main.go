package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/schemalex/git-schemalex"
)

var (
	workspace = flag.String("workspace", "", "workspace of git")
	deploy    = flag.Bool("deploy", false, "deploy")
	dsn       = flag.String("dsn", "", "")
	table     = flag.String("table", "git_schemalex_version", "table of git revision")
	schema    = flag.String("schema", "", "path to schema file")
)

func main() {
	flag.Parse()
	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-sigCh:
			cancel()
			return
		}
	}()

	r := &gitschemalex.Runner{
		Workspace: *workspace,
		Deploy:    *deploy,
		DSN:       *dsn,
		Table:     *table,
		Schema:    *schema,
	}
	err := r.Run(ctx)
	if err == gitschemalex.ErrEqualVersion {
		fmt.Println(err.Error())
		return nil
	}
	return err
}
