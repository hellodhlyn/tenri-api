package models

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hellodhlyn/tenri-api/internal/utils"
	"time"
)

type Task struct {
	UUID      string    `json:"uuid"`
	Text      string    `json:"text"`
	DueAt     time.Time `json:"dueAt"`
	CreatedAt time.Time `json:"createdAt"`
}

const taskListKey = "tasks.list"
const taskKeyPrefix = "tasks.task."

func NewTask(text string, dueAt time.Time) *Task {
	return &Task{
		UUID:      uuid.New().String(),
		Text:      text,
		DueAt:     dueAt.Round(time.Second),
		CreatedAt: time.Now().Round(time.Second),
	}
}

func GetTasks(ctx context.Context, rdb *redis.Client) ([]Task, error) {
	uuids, err := rdb.ZRange(ctx, taskListKey, 0, -1).Result()
	if err != nil {
		return nil, err
	} else if len(uuids) == 0 {
		return []Task{}, nil
	}

	keys := utils.MapSlice[string, string](uuids, func(each string) string { return taskKeyPrefix + each })
	strings, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	return utils.MapSlice[interface{}, Task](strings, func(each interface{}) Task {
		var task Task
		_ = json.Unmarshal([]byte(each.(string)), &task)
		return task
	}), nil
}

func SaveTask(ctx context.Context, rdb *redis.Client, task *Task) error {
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		bytes, err := json.Marshal(task)
		if err != nil {
			return err
		}

		pipe.Set(ctx, taskKeyPrefix+task.UUID, string(bytes), -1)
		pipe.ZAdd(ctx, taskListKey, &redis.Z{
			Score:  float64(task.DueAt.Unix()),
			Member: task.UUID,
		})
		return nil
	})
	return err
}

func DeleteTask(ctx context.Context, rdb *redis.Client, uuid string) error {
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, taskKeyPrefix+uuid)
		pipe.ZRem(ctx, taskListKey, uuid)
		return nil
	})
	return err
}
