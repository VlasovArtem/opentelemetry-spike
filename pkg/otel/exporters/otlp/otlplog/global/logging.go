package global

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"log"
	"os"
	"sync/atomic"
	"unsafe"
)

var globalLogger unsafe.Pointer

func init() {
	SetLogger(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))
}

// SetLogger overrides the globalLogger with l.
//
// To see Info messages use a logger with `l.V(1).Enabled() == true`
// To see Debug messages use a logger with `l.V(5).Enabled() == true`.
func SetLogger(l logr.Logger) {
	atomic.StorePointer(&globalLogger, unsafe.Pointer(&l))
}

func getLogger() logr.Logger {
	return *(*logr.Logger)(atomic.LoadPointer(&globalLogger))
}

// Info prints messages about the general state of the API or SDK.
// This should usually be less then 5 messages a minute.
func Info(msg string, keysAndValues ...interface{}) {
	getLogger().V(1).Info(msg, keysAndValues...)
}

// Error prints messages about exceptional states of the API or SDK.
func Error(err error, msg string, keysAndValues ...interface{}) {
	getLogger().Error(err, msg, keysAndValues...)
}

// Debug prints messages about all internal changes in the API or SDK.
func Debug(msg string, keysAndValues ...interface{}) {
	getLogger().V(5).Info(msg, keysAndValues...)
}
