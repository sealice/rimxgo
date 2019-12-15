package config

import (
	"time"

	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.AutomaticEnv()
	v.SetConfigType("yaml")

	if err := readFile("conf", 3); err != nil {
		panic(err)
	}

	if v.GetBool("Debug") {
		v.Set("LogLevel", "debug") // 调试模式始终打印Debug日志
	}
}

func readFile(path string, level uint) error {
	v.AddConfigPath(path)
	if err := v.ReadInConfig(); err != nil {
		path, level = "../"+path, level-1
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && level > 0 {
			return readFile(path, level)
		}
		return err
	}
	return nil
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return v.Get(key) }

// Sub returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func Sub(key string) *viper.Viper { return v.Sub(key) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return v.GetString(key) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return v.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return v.GetInt(key) }

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 { return v.GetInt32(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return v.GetInt64(key) }

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint { return v.GetUint(key) }

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 { return v.GetUint32(key) }

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 { return v.GetUint64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return v.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return v.GetTime(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return v.GetDuration(key) }

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int { return v.GetIntSlice(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return v.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return v.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return v.GetStringMapString(key) }

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint { return v.GetSizeInBytes(key) }

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalKey(key, rawVal, opts...)
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return v.Unmarshal(rawVal, opts...)
}

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool { return v.IsSet(key) }

func GetStringDefault(key string, value string) string {
	val := v.GetString(key)
	if val != "" {
		return val
	}
	return value
}

func GetBoolDefault(key string, value bool) bool {
	if v.IsSet(key) {
		return v.GetBool(key)
	}
	return value
}

func GetIntDefault(key string, value int) int {
	if v.IsSet(key) {
		return v.GetInt(key)
	}
	return value
}

func GetInt32Default(key string, value int32) int32 {
	if v.IsSet(key) {
		return v.GetInt32(key)
	}
	return value
}

func GetInt64Default(key string, value int64) int64 {
	if v.IsSet(key) {
		return v.GetInt64(key)
	}
	return value
}

func GetUintDefault(key string, value uint) uint {
	if v.IsSet(key) {
		return v.GetUint(key)
	}
	return value
}

func GetUint32Default(key string, value uint32) uint32 {
	if v.IsSet(key) {
		return v.GetUint32(key)
	}
	return value
}

func GetUint64Default(key string, value uint64) uint64 {
	if v.IsSet(key) {
		return v.GetUint64(key)
	}
	return value
}
