// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// iovkov@antsgames.com

package cache

import (
	"testing"
	"time"
)

func TestKeyGen(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{
			"strings in key",
			[]string{"k1", "k2", "k3"},
			"eb98f46aca624b1e402947677c54d025dc463a67",
		}, {
			"int64 in key",
			[]int64{1, 2, 3},
			"6d780b01458b623aa5f77db71ac9a02ff1d5ecda",
		}, {
			"mixed key",
			[]interface{}{1, "k2", 3},
			"ebdf5f5fd2817d5534d668548298bd135ae3dd0e",
		}, {
			"mixed key",
			[]interface{}{1, "k2", 3, map[int]string{1: "a"}},
			"60b428a15bf1e9eca224ef24e37c38ec8d8f86f9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KeyGen(tt.args)
			if k != tt.want {
				t.Errorf("KeyGen(%v), want %v, got %v ", tt.args, tt.want, k)
			}
		})
	}
}

func TestSetCacheExpiration(t *testing.T) {
	SetCacheExpiration(2 * time.Minute)
}

func TestSetAndDeleteString(t *testing.T) {
	testVal := "testString"
	key := KeyGen("test", 1)
	if err := Set(key, testVal); err != nil {
		t.Errorf("Set error:%v", err)
		return
	}
	if err := Delete(key); err != nil {
		t.Errorf("Delete error:%v", err)
		return
	}
}

func TestSetAndGetString(t *testing.T) {
	SetCacheExpiration(10 * time.Minute)
	testVal := "testString"
	key := KeyGen("test", 2)
	if err := Set(key, testVal); err != nil {
		t.Errorf("Set error:%v", err)
		return
	}
	s, _, err := GetString(key)
	if err != nil {
		t.Errorf("Get error:%v", err)
		return
	}
	if s != testVal {
		t.Error("Value not similar")
		t.Log(s, testVal)
		return
	}
}

func TestSetAndGetInt(t *testing.T) {
	tests := []int64{
		-566535465128000000,
		1329,
		3,
		323525235,
		23523052805235802,
		-323,
		-1,
	}
	for i, testVal := range tests {
		t.Run("Set And Get Int", func(t *testing.T) {
			key := KeyGen("k", i)
			if err := Set(key, testVal); err != nil {
				t.Errorf("Set error:%v", err)
				return
			}
			s, _, err := GetInt64(key)
			if err != nil {
				t.Errorf("Get error:%v", err)
				return
			}
			if s != testVal {
				t.Error("Value not similar")
				t.Log(s, testVal)
			}
		})
	}
}

func TestSetAndGetObj(t *testing.T) {
	type d struct {
		Name string
		Val  float64
	}
	testVal := d{"test", 65.5}
	key := KeyGen("test", 3)

	if err := Set(key, testVal); err != nil {
		t.Errorf("Set error:%v", err)
		return
	}

	fromC := d{}
	_, err := Load(key, &fromC)
	if err != nil {
		t.Errorf("Get error:%v", err)
		return
	}
	if fromC != testVal {
		t.Error("Value not similar")
		t.Log(fromC, testVal)
	}
}

func TestGetNotExistsKey(t *testing.T) {
	key := KeyGen("key12")
	_, ok, _ := GetInt64(key)
	if ok != false {
		t.Errorf("Key not set, but return like set")
		return
	}

}
