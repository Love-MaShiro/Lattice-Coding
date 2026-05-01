package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"
)

func GenerateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateShortID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateRunID() string {
	return "run-" + time.Now().Format("20060102150405") + "-" + GenerateShortID()
}

func GenerateSessionID() string {
	return "sess-" + GenerateShortID()
}

func CurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func GetEnv(key, defaultValue string) string {
	if value, exists := lookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func lookupEnv(key string) (string, bool) {
	return "", false
}
