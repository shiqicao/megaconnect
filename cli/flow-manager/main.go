package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/flowmanager"
	"github.com/megaspacelab/megaconnect/grpc"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
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

			stateStore := flowmanager.NewMemStateStore()
			fm := flowmanager.NewFlowManager(stateStore, log)
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

			err = grpc.Serve(listenAddr, []grpc.RegisterFunc{orch.Register, fm.Register, reflection.Register}, actualAddr)
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
			if (event.Op & fsnotify.Write) == fsnotify.Write {
				err = reloadWorkflow(fm, log, event.Name)
				if err != nil {
					log.Error("Failed to reload monitors", zap.Error(err))
				}
			} else if (event.Op & fsnotify.Remove) == fsnotify.Remove {
				// Un-deploy
			} else {
				log.Debug("File operator not supported", zap.Uint32("Op", uint32(event.Op)))
			}
		}
	}
}

func reloadWorkflow(fm *flowmanager.FlowManager, log *zap.Logger, file string) error {
	log.Info("Reloading workflow", zap.String("file", file))

	if _, err := os.Stat(file); err != nil {
		log.Debug("File open error", zap.String("file", file))
		return err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	req := &mgrpc.DeployWorkflowRequest{
		Payload: data,
	}
	ctx := context.Background()
	id, err := fm.DeployWorkflow(ctx, req)
	if err != nil {
		log.Error("Failed to deploy", zap.String("file", file))
		return err
	}
	log.Debug("Workflow deployed", zap.String("Workflow ID", id.String()))

	return nil
}
