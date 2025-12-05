package service_test

import (
	"context"
	"task-manager/internal/service"
	"task-manager/pkg/rest"
	"testing"
)

// BenchmarkCreateTask benchmarks the Create method
func BenchmarkCreateTask(b *testing.B) {
	svc := service.MakeNewTaskService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Create(context.Background(), service.RandomTask())
	}
}

// BenchmarkUpdateTask benchmarks the Update method
func BenchmarkUpdateTask(b *testing.B) {
	svc := service.MakeNewTaskService()
	task := service.CreateTestTask() // creates one test task

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.Title = "Updated Title"
		_, _ = svc.Update(context.Background(), task)
	}
}

// BenchmarkDeleteTask benchmarks the Delete method
func BenchmarkDeleteTask(b *testing.B) {
	svc := service.MakeNewTaskService()
	_ = service.CreateTestTask() // creates one test task

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// recreate task each iteration to avoid deleting same ID
		taskToDelete := service.CreateTestTask()
		_ = svc.Delete(context.Background(), taskToDelete.ID)
	}
}

// BenchmarkGetTaskByID benchmarks the GetByID method
func BenchmarkGetTaskByID(b *testing.B) {
	svc := service.MakeNewTaskService()
	task := service.CreateTestTask()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetByID(context.Background(), task.ID)
	}
}

// BenchmarkListTasks benchmarks the List method
func BenchmarkListTasks(b *testing.B) {
	svc := service.MakeNewTaskService()
	for i := 0; i < 100; i++ {
		_ = service.CreateTestTask()
	}

	query := rest.Query{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = svc.List(context.Background(), query)
	}
}

// go test -bench=. -benchmem ./internal/service
// go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
// go tool pprof http://localhost:8080/debug/pprof/heap
