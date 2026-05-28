// Package main is the temporal lesson `l3_signals_and_queries` homework scaffold for Vibe Learn.
//
// Задача: OrderWorkflow с состоянием items: сигнал add-item, query current-items, Update с валидатором.
// Реализуй workflow и активности ниже — сигнатуры и тестовая поверхность
// фиксированы; CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// SDK: go.temporal.io/sdk (worker + workflow + activity).
// Воркер подключается к Temporal по TEMPORAL_ADDRESS (дефолт localhost:7233 —
// совпадает с docker-compose.yml) и слушает task queue из TaskQueue().
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// TemporalAddress — адрес Temporal frontend. Дефолт совпадает с docker-compose.yml.
func TemporalAddress() string {
	return envOr("TEMPORAL_ADDRESS", "localhost:7233")
}

// TaskQueue — очередь задач, которую слушает воркер этого урока.
func TaskQueue() string {
	return envOr("TEMPORAL_TASK_QUEUE", "lesson-l3_signals_and_queries-tq")
}

// ----- Workflow: OrderWorkflow -----
//
// Оркеструет активности ниже. Тело — TODO: добавь ExecuteActivity-шаги,
// ActivityOptions (StartToCloseTimeout, RetryPolicy) и обработку ошибок
// согласно README.md. Должно оставаться ДЕТЕРМИНИРОВАННЫМ (никаких
// time.Now/rand/итераций по map — используй workflow.Now/SideEffect).
func OrderWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("OrderWorkflow started", "taskQueue", TaskQueue())

	// TODO #1: вызови активность CheckStock через workflow.ExecuteActivity.
	// var checkstockRes string
	// if err := workflow.ExecuteActivity(ctx, CheckStock).Get(ctx, &checkstockRes); err != nil {
	// 	return err
	// }
	// TODO #2: вызови активность Checkout через workflow.ExecuteActivity.
	// var checkoutRes string
	// if err := workflow.ExecuteActivity(ctx, Checkout).Get(ctx, &checkoutRes); err != nil {
	// 	return err
	// }

	return nil
}

// ----- Activity #1: CheckStock -----
//
// проверить наличие товара на складе (валидатор Update отклоняет отсутствующий)
func CheckStock(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("CheckStock: not implemented")
}

// ----- Activity #2: Checkout -----
//
// оформить заказ по накопленным items
func Checkout(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("Checkout: not implemented")
}

// ----- main entry: register worker + run with graceful shutdown -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — temporal lesson %s scaffold up", "l3_signals_and_queries")
	log.Printf("temporal address: %s  task queue: %s", TemporalAddress(), TaskQueue())
	log.Printf("Реализуй workflow и активности, затем `go test ./...`. README.md содержит задачу.")

	c, err := client.Dial(client.Options{HostPort: TemporalAddress()})
	if err != nil {
		log.Fatalf("unable to create Temporal client (is `docker compose up -d` running?): %v", err)
	}
	defer c.Close()

	w := worker.New(c, TaskQueue(), worker.Options{})
	w.RegisterWorkflow(OrderWorkflow)
	w.RegisterActivity(CheckStock)
	w.RegisterActivity(Checkout)

	// Graceful shutdown so `go run .` is interactive — worker.InterruptCh()
	// stops the worker on Ctrl-C / SIGTERM.
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped with error: %v", err)
	}
}
