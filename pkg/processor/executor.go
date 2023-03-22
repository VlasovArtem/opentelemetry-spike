package main

import (
	"bytes"
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os/exec"
	"spike-go-opentelemetry-logging/pkg/common"
	"strconv"
)

var executorLog = otelzap.New(zap.NewExample(),
	otelzap.WithMinLevel(zap.InfoLevel),
	otelzap.WithTraceIDField(true),
	otelzap.WithCaller(false),
)

func execute(ctx context.Context, request insertDataRequest) {
	executorCtxt, span := tracer.Start(ctx, "execute", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	otelzap.Ctx(executorCtxt).Info("Executing command", zap.Int("random", request.Random))
	var logs bytes.Buffer

	cmd := exec.Command(common.GlobalOpts.Executor.BasePath+"/container_execution.sh", strconv.Itoa(request.Random))
	cmd.Stderr = &logs
	cmd.Stdout = &logs

	err := cmd.Start()
	if err != nil {
		executorLog.Ctx(executorCtxt).Error("Error starting command", zap.Error(err))
	} else {
		err = cmd.Wait()
		if err != nil {
			executorLog.Ctx(executorCtxt).Error("Error waiting for command", zap.Error(err))
			executorLog.Ctx(executorCtxt).Error("Logs", zap.ByteString("logs", logs.Bytes()))
		} else {
			executorLog.Ctx(executorCtxt).Info("Logs", zap.ByteString("logs", logs.Bytes()))
		}
	}
	executorLog.Ctx(executorCtxt).Info("Command executed")
}
