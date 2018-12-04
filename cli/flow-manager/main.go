package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"time"

	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/flowmanager"
	"github.com/megaspacelab/megaconnect/grpc"

	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	listenPortFlag = cli.IntFlag{
		Name:  "listen-port",
		Usage: "Listening port",
		Value: 9000,
	}
	mblockIntervalFlag = cli.DurationFlag{
		Name:  "mblock-interval",
		Usage: "MBlock interval",
		Value: 10 * time.Second,
	}
)

func main() {
	app := cli.App{
		Usage: "Megaspace flow manager",
		Flags: []cli.Flag{
			&mcli.DebugFlag,
			&mcli.DataDirFlag,
			&listenPortFlag,
			&mblockIntervalFlag,
		},
		Action: mcli.ToExitCode(run),
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	log, err := mcli.NewLogger(ctx.Bool(mcli.DebugFlag.Name))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	dataDir := ctx.String(mcli.DataDirFlag.Name)
	stateStore, err := flowmanager.NewSQLiteStateStore(path.Join(dataDir, "fm", "state", "state.db"))
	if err != nil {
		log.Error("Failed to create state store", zap.Error(err))
		return err
	}
	defer stateStore.Close()

	fm := flowmanager.NewFlowManager(stateStore, log)
	orch := flowmanager.NewOrchestrator(fm, log)

	mblockTicker := time.NewTicker(ctx.Duration(mblockIntervalFlag.Name))
	defer mblockTicker.Stop()

	go func() {
		for range mblockTicker.C {
			mblock, err := fm.FinalizeAndCommitMBlock()
			if err != nil {
				log.Error("Failed to finalize mblock", zap.Error(err))
				continue
			}

			log.Info("Finalized mblock", zap.Int64("height", mblock.Height), zap.Int("numEvents", len(mblock.Events)))
		}
	}()

	listenAddr := fmt.Sprintf(":%d", ctx.Int(listenPortFlag.Name))
	actualAddr := make(chan net.Addr, 1)
	go func() {
		addr, ok := <-actualAddr
		if ok {
			log.Debug("Serving", zap.Stringer("listenAddr", addr))
		}
	}()

	err = grpc.Serve(listenAddr, []grpc.RegisterFunc{orch.Register, fm.Register, reflection.Register}, actualAddr)
	if err != nil {
		log.Error("Failed to serve", zap.Error(err))
		return err
	}

	return nil
}
