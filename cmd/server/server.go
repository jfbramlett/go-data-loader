package main

import (
	"context"
	"fmt"
	"github.com/jfbramlett/go-aop/pkg/aop"
	"github.com/jfbramlett/go-aop/pkg/logging"
	"github.com/jfbramlett/go-aop/pkg/metrics"
	"github.com/jfbramlett/go-aop/pkg/tracing"
	server "github.com/jfbramlett/go-template/pkg/api"
	"github.com/spf13/cobra"
)

// The serveCmd will execute the generate command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the web server",
	Run:   serve,
}

func init() {
	// Here we create the command line flags for our app, and bind them to our package-local
	// config variable.
	flags := serveCmd.Flags()
	flags.String("service_name", "gotemplate", "name of the service")
	flags.Bool("metrics_enabled", true, "if metrics should be produced")
	flags.Float64("trace_sample_rate", 1.0, "the rate to sample trace at (percentage as decimal between 0 and 1)")
	flags.String("web_root", "./", "root directory for serving our content")

	// Add the "serve" sub-command to the root command.
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	logging.InitLogging()
	logger := logging.GetLogger(context.Background())

	service, err := cmd.Flags().GetString("service_name")
	if err != nil {
		logger.Error(err, "failed to resolve service name")
		return
	}

	webRoot, err := cmd.Flags().GetString("web_root")
	if err != nil {
		logger.Error(err, "failed to resolve web root")
		return
	}

	aop.InitAOP(service)
	aop.RegisterJoinPoint(aop.NewRegexPointcut(".*"), aop.NewSpanFuncAdvice())
	aop.RegisterJoinPoint(aop.NewRegexPointcut(".*"), aop.NewTimedFuncAdvice(service, fmt.Sprintf("%s metrics", service)))
	aop.RegisterJoinPoint(aop.NewRegexPointcut(".*"), aop.NewLoggingFuncAdvice())

	metrics.InitMetrics(metrics.MetricsConfig{Port: 1000, URL: "/metrics"})
	tracing.InitTracing(tracing.TracingConfig{Service: service, ReporterPort: 8090, ReporterUrl: "http://localhost:9411/api/v2/spans"})

	runner := server.NewWebRunner(webRoot)
	runner.Run()
}
