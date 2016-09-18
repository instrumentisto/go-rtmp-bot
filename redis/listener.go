package redis

import (
	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"gopkg.in/redis.v4"
	"log"
)

const (
	START_COMMAND       = "start_test"         // Stress test started.
	STOP_COMMAND        = "stop_test"          // Stress test stopped.
	STRESS_TEST_CHANNEL = "stress_test_client" // Redis channel name.
)

// Redis db pub/sub listener.
type RedisListener struct {
	client      *redis.Client // Redis client
	app_handler controller.AppHandler
}

// Returns new instance of Redis listener.
// props: r_url    string   Redis server URL.
//        password string   Redis server password.
//        db       int      Redis server Data Base ID.
func NewRedisListener(
	r_url string,
	password string,
	db int,
	handler controller.AppHandler) *RedisListener {
	return &RedisListener{
		app_handler: handler,
		client: redis.NewClient(&redis.Options{
			Addr:     r_url,
			Password: password,
			DB:       db,
		}),
	}
}

// Listens Redis pub/sub channel.
func (l *RedisListener) Listen() {
	defer l.client.Close()
	pubsub, err := l.client.Subscribe(STRESS_TEST_CHANNEL)
	if err != nil {
		log.Printf("CLIENT SUBSCRIBE ERROR: %v", err)
		return
	}
	defer pubsub.Unsubscribe(STRESS_TEST_CHANNEL)
	for {
		err := l.ping()
		if err != nil {
			log.Printf("CLIENT PING ERROR: %s", err.Error())
			return
		}
		mess, err := pubsub.ReceiveMessage()
		if err == nil {
			l.readRedisMessage(mess)
		}
	}
}

// Calls redis pub/sub channel.
//
// param: command string   Any command for calls.
func (l *RedisListener) Call(command string) error {
	return l.client.Publish(STRESS_TEST_CHANNEL, command).Err()
}

// Writes to redis db any value by key.
//
// params: key   string   Key of value.
//         value int      Integer value.
func (l *RedisListener) Write(key string, value int) error {
	return l.client.Set(key, value, 0).Err()
}

// Writes to redis map with name any value by key.
//
// params: map_name   string   Name of the redis hash map.
//         field_name string   Name of field.
//         value      string   Any string value.
func (l *RedisListener) WriteToMap(
	map_name string, field_name string, value string) error {
	return l.client.HSet(map_name, field_name, value).Err()
}

// Returns map from redis.
//
// params: map_name string   Redis map name.
func (l *RedisListener) GetMap(map_name string) (map[string]string, error) {
	return l.client.HGetAll(map_name).Result()
}

// Reads integer value from redis.
//
// params: key string   The values key.
func (l *RedisListener) Read(key string) (int64, error) {
	return l.client.Get(key).Int64()
}

// Closes listener
func (l *RedisListener) Close() {
	l.client.Close()
}

// Reads Redis pub/sub messages.
// Just writes new signal to application signal handler.
func (l *RedisListener) readRedisMessage(mess *redis.Message) {
	signal := model.NewSignal(mess.Payload, "redis")
	l.app_handler.OnSignal(signal)
}

// Pings redis pub/sub channel.
func (l *RedisListener) ping() error {
	err := l.client.Ping().Err()
	if err != nil {
		return err
	}
	return nil
}
