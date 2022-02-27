package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type Logger interface {
	Log(critical bool, message string)
}

type TrunkLogger struct {
	address   string
	ctx       context.Context
	backup    *os.File
	active    bool
	prefixLen int
}

func TrunkLoggerInit(address string, ctx context.Context, backup *os.File, prefixLen int) *TrunkLogger {
	if prefixLen == 0 {
		log.Fatal("(TrunkLoggerInit): Prefix length must be greater than 0!")
	}

	return &TrunkLogger{
		address,
		ctx,
		backup,
		true,
		prefixLen,
	}
}

func (l *TrunkLogger) Log(critical bool, message string) {

	// Guard case if trunk's logging service
	// failed to respond on a previous attempt
	if !l.active {
		_, _ = l.backup.Write([]byte(message))
		return
	}

	stream, err := net.Dial("tcp", l.address)
	if err != nil {
		l.active = false                // Haven't found docs that describe what cases it returns err
		l.backup.Write([]byte(message)) // So set trunk to 'inactive' to be safe
		return
	}

	// Considering adding stream to struct, depends on trunk's conn lifespan
	defer stream.Close()

	req := fmtReq(critical, message)

	// Use logger's given context. Open to lowering the timeout, 1
	// second might be unnecessarily long for just a single logging operation.
	ctx, cancel := context.WithTimeout(l.ctx, time.Second)
	defer cancel()

	res := make(chan bool)
	defer close(res)

	// TODO: Log debugging information with context to stdout?
	go func() {
		_, err = stream.Write(req)
		if err != nil {
			res <- false
			return
		}

		prefixBuff := make([]byte, l.prefixLen)

		_, err = stream.Read(prefixBuff)
		if err != nil && err != io.EOF {
			res <- false
			return
		}

		msgLen, err := strconv.Atoi(string(prefixBuff))
		if err != nil || msgLen == 0 {
			res <- false
			return
		}

		msgBuff := make([]byte, msgLen)

		_, err = stream.Read(msgBuff)
		if err != nil && err != io.EOF {
			res <- false
			return
		}

		// Full message if successful should be "00Ok" but "00" is all we need
		// to confirm success, and the "Ok" might change later, while the "00" won't.
		if string(msgBuff[0:2]) == "00" {
			res <- true
			return
		} else {
			res <- false
			return
		}
	}()

	for {
		select {
		case dst := <-res:
			if !dst {
				l.active = false
				_, _ = l.backup.Write([]byte(message))
			}

			return
		case <-ctx.Done():
			cancel()

			// If there ends up being high variance in trunk's response time,
			// response time might go over 1 second,
			// so setting trunk inactive here may be dumping it prematurely.
			// If trunk starts dropping off in response time the longer it runs,
			// this service doesn't take a hit anyway, so don't set to inactive.
			// l.active = false
			_, _ = l.backup.Write([]byte(message))
			return
		}
	}
}

func fmtReq(critical bool, message string) []byte {
	var level string
	if critical {
		level = "!"
	} else {
		level = "$"
	}

	return []byte("lo" + level + message + "\n")
}
