/*
 * © 2023 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/snyk/snyk-ls/internal/types"
)

type lspWriter struct {
	writeChan chan types.LogMessageParams
	readyChan chan bool
	server    types.Server
}

func New(server types.Server) zerolog.LevelWriter {
	if server != nil {
		_, _ = fmt.Fprintln(os.Stderr, "LSP logger: starting with non-nil server")
	}
	readyChan := make(chan bool)
	writeChan := make(chan types.LogMessageParams, 1000000)
	w := &lspWriter{
		writeChan: writeChan,
		readyChan: readyChan,
		server:    server,
	}
	go w.startServerSenderRoutine()
	// let the routine startup first
	_, _ = fmt.Fprintln(os.Stderr, "LSP logger: waiting for ready signal...")
	<-w.readyChan
	_, _ = fmt.Fprintln(os.Stderr, "LSP logger: started")
	return w
}

func (w *lspWriter) Write(p []byte) (n int, err error) {
	return w.WriteLevel(zerolog.InfoLevel, p)
}

func (w *lspWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	levelEnabled := level > zerolog.TraceLevel && level < zerolog.NoLevel
	if w.server != nil && levelEnabled {
		w.writeChan <- types.LogMessageParams{
			Type:    mapLogLevel(level),
			Message: string(p),
		}
		return len(p), nil
	}

	if levelEnabled {
		return os.Stderr.Write(p)
	}

	return 0, nil
}

func (w *lspWriter) startServerSenderRoutine() {
	w.readyChan <- true
	var err error
	for msg := range w.writeChan {
		err = w.server.Notify(context.Background(), "window/logMessage", msg)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(msg.Message))
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "LSP logger (%p) stopped\n", w)
}

func mapLogLevel(level zerolog.Level) (mt types.MessageType) {
	switch level {
	case zerolog.PanicLevel:
		fallthrough
	case zerolog.FatalLevel:
		fallthrough
	case zerolog.ErrorLevel:
		mt = types.Error
	case zerolog.WarnLevel:
		mt = types.Warning
	case zerolog.InfoLevel:
		mt = types.Info
	case zerolog.DebugLevel:
		mt = types.Log
	default:
		mt = 0
	}
	return mt
}
