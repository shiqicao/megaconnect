package main

import (
	"fmt"
	"net"
	"os"

	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/flowmanager"
	"github.com/megaspacelab/megaconnect/grpc"
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
		},
	}

	app.Run(os.Args)
}
