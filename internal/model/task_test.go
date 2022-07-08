package model_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/hellodhlyn/tenri-api/internal/model"
	"github.com/hellodhlyn/tenri-api/internal/utils"
	"time"

	"testing"
)

func TestGetTasks(t *testing.T) {
	ctx := context.TODO()
	rdb, mock := redismock.NewClientMock()
	uuids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	mock.ExpectZRange("tasks.list", 0, -1).SetVal(uuids)
	mock.ExpectMGet(utils.MapSlice[string, string](uuids, func(each string) string { return "tasks.task." + each })...).
		SetVal(utils.MapSlice[string, interface{}](uuids, func(each string) interface{} {
			bytes, _ := json.Marshal(&model.Task{
				UUID: each, Text: faker.Sentence(), DueAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now(),
			})
			return string(bytes)
		}))

	results, err := model.GetTasks(ctx, rdb)
	if err != nil {
		t.Error(err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	for i, result := range results {
		if uuids[i] != result.UUID {
			t.Errorf("uuid not matched (expected: %s, actual: %s)", uuids[i], result.UUID)
		}
	}
}

func TestSaveTask_Create(t *testing.T) {
	ctx := context.TODO()
	rdb, mock := redismock.NewClientMock()

	now := time.Now().Round(time.Second)
	task := model.Task{
		UUID:      uuid.New().String(),
		Text:      faker.Sentence(),
		DueAt:     now.Add(24 * time.Hour),
		CreatedAt: now,
	}

	expectedJSON := fmt.Sprintf(
		`{"uuid":"%s","text":"%s","dueAt":"%s","createdAt":"%s"}`,
		task.UUID,
		task.Text,
		task.DueAt.Format(time.RFC3339),
		task.CreatedAt.Format(time.RFC3339),
	)
	mock.ExpectTxPipeline()
	mock.ExpectSet("tasks.task."+task.UUID, expectedJSON, -1).SetVal("")
	mock.ExpectZAdd("tasks.list", &redis.Z{Score: float64(task.DueAt.Unix()), Member: task.UUID}).SetVal(1)
	mock.ExpectTxPipelineExec()

	err := model.SaveTask(ctx, rdb, &task)
	if err != nil {
		t.Error(err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
