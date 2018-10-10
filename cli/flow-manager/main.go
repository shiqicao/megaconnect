package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/megaspacelab/megaconnect/flowmanager"
	"github.com/megaspacelab/megaconnect/grpc"
	wf "github.com/megaspacelab/megaconnect/workflow"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

const (
	listenAddr = ":12345"
)

// TODO - make this a proper cli app.
func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger. %v", err)
		return
	}

	flowManager := flowmanager.NewFlowManager(log)
	orchestrator := flowmanager.NewOrchestrator(flowManager, log)

	done := make(chan struct{}, 1)
	go loadAndWatchChainConfigs(flowManager, log, done, "Example", "Ethereum", "Bitcoin")

	log.Debug("Serving", zap.String("listenAddr", listenAddr))
	err = grpc.Serve(listenAddr, []grpc.RegisterFunc{orchestrator.Register, reflection.Register}, nil)
	done <- struct{}{}
	if err != nil {
		log.Error("Failed to serve", zap.Error(err))
	}
}

func loadAndWatchChainConfigs(
	flowManager *flowmanager.FlowManager,
	log *zap.Logger,
	done <-chan struct{},
	chains ...string,
) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	for _, chain := range chains {
		chainDir := path.Join(usr.HomeDir, ".megaspace", "fm", "chains", chain)
		err = os.MkdirAll(chainDir, 0700)
		if err != nil {
			panic(err)
		}

		err = reloadMonitors(flowManager, log, chain, path.Join(chainDir, "monitors"))
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
			err = reloadMonitors(flowManager, log, chain, event.Name)
			if err != nil {
				log.Error("Failed to reload monitors", zap.String("chain", chain))
			}
		}
	}
}

func reloadMonitors(flowManager *flowmanager.FlowManager, log *zap.Logger, chain, file string) error {
	log.Info("Reloading monitors", zap.String("chain", chain), zap.String("file", file))

	fs, err := os.Open(file)
	defer fs.Close()

	if _, ok := err.(*os.PathError); ok {
		log.Debug("File does not exist", zap.String("file", file))
		flowManager.SetChainConfig(chain, nil, nil)
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

	flowManager.SetChainConfig(chain, monitors, nil)
	return nil
}
