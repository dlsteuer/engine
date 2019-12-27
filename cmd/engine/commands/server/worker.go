package server

import (
	"context"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/battlesnakeio/engine/rules"

	"github.com/battlesnakeio/engine/controller/pb"
	"github.com/battlesnakeio/engine/worker"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	promgrpc "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const pingRetryDelay = 1 * time.Second

var (
	workerThreads      = 10
	workerPollInterval = 1 * time.Second
	workerChaos        = false
)

func init() {
	workerCmd.Flags().IntVarP(&workerThreads, "threads", "t", workerThreads, "worker processor threads, this is the amount of concurrent games a worker can process")
	workerCmd.Flags().StringVarP(&controllerAddr, "controller-addr", "c", controllerAddr, "address of the controller")
	workerCmd.Flags().DurationVarP(&workerPollInterval, "poll-interval", "p", workerPollInterval, "worker poll interval")
	workerCmd.Flags().BoolVar(&workerChaos, "chaos", workerChaos, "introduce chaotic latency into the worker")
	RootCmd.Flags().AddFlagSet(workerCmd.Flags())
}

// randTimeoutInterceptor provides a random amount of variance to all GRPC calls
// at the client level. This is part of the chaos mode for the workers. It means
// that calls will randomly go over the lock interval triggering some
// interesting situations.
func randTimeoutInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var sleep time.Duration
	if rand.Intn(100) <= 20 {
		sleep = time.Duration(rand.Intn(5)) * time.Second
	} else {
		sleep = time.Duration(rand.Intn(50)) * time.Millisecond
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleep):
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

var workerCmd = &cobra.Command{
	Use:    "worker",
	Short:  "runs the engine worker",
	PreRun: func(c *cobra.Command, args []string) { prometheus() },
	Run: func(c *cobra.Command, args []string) {
		interceptors := []grpc.UnaryClientInterceptor{promgrpc.UnaryClientInterceptor}
		if workerChaos {
			log.Warn("using chaos mode")
			interceptors = append(interceptors, randTimeoutInterceptor)
		}
		client, err := pb.Dial(controllerAddr, grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(interceptors...),
		))
		if err != nil {
			log.WithError(err).
				WithField("address", controllerAddr).
				Fatal("failed to dial controller")
		}

		w := &worker.Worker{
			ControllerClient: client,
			PollInterval:     workerPollInterval,
			Rulesets:         initializeRulesets(),
		}
		w.RunGame = w.Runner

		// Begin pinging controller to push useful logs to an operator.
		go func() {
			for {
				ping, err := client.Ping(context.Background(), &pb.PingRequest{})
				if err == nil {
					log.WithField("version", ping.Version).
						Info("connection to controller healthy")
					break
				} else {
					log.WithError(err).Warn("controller connection unhealthy")
					time.Sleep(pingRetryDelay)
				}
			}
		}()

		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(workerThreads)

		for i := 0; i < workerThreads; i++ {
			go func(i int) {
				log.WithField("worker", i).Info("Battlesnake worker starting")
				w.Run(ctx, i)
				wg.Done()
			}(i)
		}
		wg.Wait()
	},
}

func initializeRulesets() map[string]rules.Ruleset {
	rulesets := map[string]rules.Ruleset{
		"standard": &rules.DefaultRuleset{},
	}

	// plugins only work on windows, so just return the default ruleset
	if runtime.GOOS == "windows" {
		return rulesets
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Error("unable to retrieve user home dir")
		return rulesets
	}
	plugins, err := filepath.Glob(path.Join(homeDir, ".battlesnake/rulesets/*.so"))
	if err != nil {
		log.WithError(err).Error("unable to retrieve plugin listing")
		return rulesets
	}
	for _, p := range plugins {
		name, rs, err := loadPlugin(p)
		if err != nil {
			log.WithError(err).WithField("plugin", p).Error("unable to load plugin")
			continue
		}

		log.WithField("name", name).Info("loaded ruleset")
		rulesets[name] = rs
	}

	return rulesets
}

func loadPlugin(pluginPath string) (name string, rs rules.Ruleset, err error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.WithError(err).Error("unable to load ruleset")
		return "", nil, err
	}

	srs, err := p.Lookup("Ruleset")
	if err != nil {
		log.WithError(err).Error("unable to find ruleset symbol in plugin")
		return "", nil, err
	}

	rs, ok := srs.(rules.Ruleset)
	if !ok {
		log.WithError(err).Error("ruleset does not match ruleset interface")
		return "", nil, err
	}

	name = strings.TrimRight(path.Base(pluginPath), ".so")
	return name, rs, nil
}
