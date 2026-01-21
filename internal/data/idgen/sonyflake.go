package idgen

import (
	"strconv"
	"time"

	"github.com/sony/sonyflake/v2"
)

type IDGenerator struct {
	idgen *sonyflake.Sonyflake
}

func NewIDGenerator() *IDGenerator {
	idgen, err := sonyflake.New(sonyflake.Settings{})
	if err != nil {
		panic(err)
	}
	return &IDGenerator{
		idgen: idgen,
	}
}

func (ig *IDGenerator) NextID(prefix string) string {
	id, err := ig.idgen.NextID()
	if err != nil {
		panic(err)
	}
	return prefix + strconv.FormatInt(id, 10)
}

func (ig *IDGenerator) NextOrderID() string {
	now := time.Now()
	prefix := now.Format("20060102150405")
	return prefix + ig.NextID("")
}
