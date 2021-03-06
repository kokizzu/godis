package godis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedis_Eval(t *testing.T) {
	initDb()
	redis := NewRedis(option)
	defer redis.Close()
	s, err := redis.Eval(`return redis.call("get",KEYS[1])`, 1, "godis")
	assert.Nil(t, err)
	assert.Equal(t, "good", s)

	s, err = redis.Eval(`return redis.call("set",KEYS[1],ARGV[1])`, 1, "eval", "godis")
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	s, err = redis.Eval(`return redis.call("get",KEYS[1])`, 1, "eval")
	assert.Nil(t, err)
	assert.Equal(t, "godis", s)

	redisBroken := NewRedis(option)
	defer redisBroken.Close()
	redisBroken.client.connection.host = "localhost1"
	redisBroken.Close()
	_, err = redisBroken.Eval(`return redis.call("get",KEYS[1])`, 1, "godis")
	assert.NotNil(t, err)
}

func TestRedis_EvalByKeyArgs(t *testing.T) {
	initDb()
	redis := NewRedis(option)
	defer redis.Close()
	s, err := redis.EvalByKeyArgs(`return redis.call("get",KEYS[1])`, []string{"godis"}, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "good", s)

	s, err = redis.EvalByKeyArgs(`return redis.call("set",KEYS[1],ARGV[1])`, []string{"eval"}, []string{"godis"})
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	s, err = redis.EvalByKeyArgs(`return redis.call("get",KEYS[1])`, []string{"eval"}, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "godis", s)
	TestRedis_Set(t)

	redisBroken := NewRedis(option)
	defer redisBroken.Close()
	redisBroken.client.connection.host = "localhost1"
	redisBroken.Close()
	_, err = redisBroken.EvalByKeyArgs(`return redis.call("get",KEYS[1])`, []string{"godis"}, []string{})
	assert.NotNil(t, err)
}

func TestRedis_ScriptLoad(t *testing.T) {
	initDb()
	redis := NewRedis(option)
	defer redis.Close()
	sha, err := redis.ScriptLoad(`return redis.call("get",KEYS[1])`)
	assert.Nil(t, err)

	bools, err := redis.ScriptExists(sha)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true}, bools)

	s, err := redis.EvalSha(sha, 1, "godis")
	assert.Nil(t, err)
	assert.Equal(t, "good", s)

	redisBroken := NewRedis(option)
	defer redisBroken.Close()
	redisBroken.client.connection.host = "localhost1"
	redisBroken.Close()
	_, err = redisBroken.ScriptLoad(`return redis.call("get",KEYS[1])`)
	assert.NotNil(t, err)
	_, err = redisBroken.ScriptExists(sha)
	assert.NotNil(t, err)
}
