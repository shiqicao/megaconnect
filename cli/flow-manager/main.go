package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"

	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/flowmanager"
	"github.com/megaspacelab/megaconnect/grpc"
	wf "github.com/megaspacelab/megaconnect/workflow"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	listenPortFlag := cli.IntFlag{
		Name:  "listen-port",
		Usage: "Listening port",
		Value: 9000,
	}

	app := cli.App{
		Usage: "Megaspace flow manager",
		Flags: []cli.Flag{
			&mcli.DebugFlag,
			&mcli.DataDirFlag,
			&listenPortFlag,
		},
		Action: func(ctx *cli.Context) error {
			log, err := mcli.NewLogger(ctx.Bool(mcli.DebugFlag.Name))
			if err != nil {
				return err
			}

			fm := flowmanager.NewFlowManager(log)
			fm.SetChainConfig("Example", nil, nil)
			fm.SetChainConfig("Ethereum", nil, nil)
			fm.SetChainConfig("Bitcoin", nil, nil)
			orch := flowmanager.NewOrchestrator(fm, log)

			dataDir := ctx.Path(mcli.DataDirFlag.Name)
			done := make(chan struct{}, 1)
			go loadAndWatchChainConfigs(fm, log, dataDir, done)

			listenAddr := fmt.Sprintf(":%d", ctx.Int(listenPortFlag.Name))
			actualAddr := make(chan net.Addr, 1)
			go func() {
				addr, ok := <-actualAddr
				if ok {
					log.Debug("Serving", zap.Stringer("listenAddr", addr))
				}
			}()

			err = grpc.Serve(listenAddr, []grpc.RegisterFunc{orch.Register, reflection.Register}, actualAddr)
			done <- struct{}{}
			if err != nil {
				log.Error("Failed to serve", zap.Error(err))
				return err
			}

			return nil
		},
	}

	app.Run(os.Args)
}

func loadAndWatchChainConfigs(
	fm *flowmanager.FlowManager,
	log *zap.Logger,
	dataDir string,
	done <-chan struct{},
) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	wfDir := path.Join(dataDir, "fm", "workflows")
	err = os.MkdirAll(wfDir, 0700)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0)
	filepath.Walk(wfDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		err := reloadWorkflow(fm, log, file)
		if err != nil {
			panic(err)
		}
	}

	err = watcher.Add(wfDir)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-done:
			return
		case event := <-watcher.Events:
			log.Info("Received fs event", zap.Stringer("event", event))
			if event.Op == fsnotify.Create {
				err = reloadWorkflow(fm, log, event.Name)
				if err != nil {
					log.Error("Failed to reload monitors", zap.Error(err))
				}
			} else if event.Op == fsnotify.Remove {
				// Un-deploy
			} else {
				log.Debug("File operator not supported", zap.Uint32("Op", uint32(event.Op)))
			}
		}
	}
}

func reloadWorkflow(fm *flowmanager.FlowManager, log *zap.Logger, file string) error {
	log.Info("Reloading workflow", zap.String("file", file))

	fs, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fs.Close()

	if _, ok := err.(*os.PathError); ok {
		log.Debug("File does not exist", zap.String("file", file))
		return nil
	} else if err != nil {
		return err
	}

	workflow, err := wf.NewDecoder(fs).DecodeWorkflow()
	if err != nil {
		return err
	}

	if workflow.Name() != path.Base(file) {
		log.Error("Workflow/file name mismatch", zap.String("workflow", workflow.Name()), zap.String("file", file))
		return fmt.Errorf("Workflow name %s mismatches file name %s", workflow.Name(), file)
	}
	if err = fm.DeployWorkflow(workflow); err != nil {
		log.Error("Failed to deploy", zap.String("workflow", workflow.Name()), zap.Error(err))
		return err
	}

	return nil
}
