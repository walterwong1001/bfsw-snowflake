package snowflake

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
// configured machine id - 10 bits - gives us up to 1024 machines
// sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)
const (
	machine_id_bits = 10                           // bit length of machine id
	sequence_bits   = 12                           // bit length of sequence number
	max_machine_id  = -1 ^ (-1 << machine_id_bits) // max machine id (1023)
	sequence_mask   = -1 ^ (-1 << sequence_bits)   // sequence mask 4095
	epoch           = 1721001600000                // 2024-07-15 00:00:00 UTC in milliseconds
)

type snowflake struct {
	mutex         sync.Mutex
	machineId     int
	sequence      int
	lastTimestamp int64
}

func newSnowflake(machineId int) (*snowflake, error) {
	if machineId > max_machine_id {
		return nil, errors.New(fmt.Sprintf("Machine Id can't be greater than %d", max_machine_id))
	}
	return &snowflake{
		machineId:     machineId,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

func (s *snowflake) NextID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixMilli() // 转换为毫秒
	if timestamp < s.lastTimestamp {
		log.Printf("clock is moving backwards.  Rejecting requests until %d.\n", s.lastTimestamp)
		return 0
	}

	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequence_mask
		if s.sequence == 0 {
			timestamp = tilNextMillis(timestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return (timestamp-epoch)<<(machine_id_bits+sequence_bits) | int64(s.machineId<<sequence_bits) | int64(s.sequence)
}

func tilNextMillis(timestamp int64) int64 {
	var current = time.Now().UnixMilli()
	for current <= timestamp {
		current = time.Now().UnixMilli()
	}
	return current
}
