package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path"

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
			orch := flowmanager.NewOrchestrator(fm, log)

			dataDir := ctx.Path(mcli.DataDirFlag.Name)
			done := make(chan struct{}, 1)
			go loadAndWatchChainConfigs(fm, log, dataDir, done, "Example", "Ethereum", "Bitcoin")

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
	chains ...string,
) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	for _, chain := range chains {
		chainDir := path.Join(dataDir, "fm", "chains", chain)
		err = os.MkdirAll(chainDir, 0700)
		if err != nil {
			panic(err)
		}

		err = reloadMonitors(fm, log, chain, path.Join(chainDir, "monitors"))
		if err != nil {
			panic(err)
		}

		err = watcher.Add(chainDir)
		if err != nil {
			panic(err)
		}
	}

	for {
		select {
		case <-done:
			return
		case event := <-watcher.Events:
			log.Info("Received fs event", zap.Stringer("event", event))
			if path.Base(event.Name) != "monitors" {
				continue
			}
			chain := path.Base(path.Dir(event.Name))
			err = reloadMonitors(fm, log, chain, event.Name)
			if err != nil {
				log.Error("Failed to reload monitors", zap.String("chain", chain))
			}
		}
	}
}

func reloadMonitors(fm *flowmanager.FlowManager, log *zap.Logger, chain, file string) error {
	log.Info("Reloading monitors", zap.String("chain", chain), zap.String("file", file))

	fs, err := os.Open(file)
	defer fs.Close()

	if _, ok := err.(*os.PathError); ok {
		log.Debug("File does not exist", zap.String("file", file))
		fm.SetChainConfig(chain, nil, nil)
		return nil
	} else if err != nil {
		return err
	}

	scanner := bufio.NewScanner(fs)
	monitorDecls := []*wf.MonitorDecl{}
	for scanner.Scan() {
		monitorRaw, err := hex.DecodeString(scanner.Text())
		if err != nil {
			return err
		}
		monitor, err := wf.NewByteDecoder(monitorRaw).DecodeMonitorDecl()
		if err != nil {
			log.Error("Failed to decode monitor", zap.Error(err))
			return err
		}
		log.Debug("Adding monitor valuations", zap.Stringer("monitor", monitor))
		monitorDecls = append(monitorDecls, monitor)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	monitors := make(map[flowmanager.MonitorID]*grpc.Monitor)
	for i, m := range monitorDecls {
		cond, err := wf.EncodeExpr(m.Condition())
		if err != nil {
			log.Error("Failed to encode condition", zap.Stringer("condition", m.Condition()))
			return err
		}
		monitors[flowmanager.MonitorID(i)] = &grpc.Monitor{
			Id: int64(i),
			// TODO: determine required expressions at chain manager
			Evaluations: [][]byte{},
			Condition:   cond,
		}
	}

	fm.SetChainConfig(chain, monitors, nil)
	return nil
}
