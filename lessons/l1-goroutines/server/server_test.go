package main

import (
	"errors"
	"io"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/ronna-s/golab2026/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServe(t *testing.T) {
	t.Run("listener returns an error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			exampleErr := errors.New("error")
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return((net.Conn)(nil), exampleErr).Once()
			assert.ErrorIs(t, exampleErr, Serve(listener, nil))
			synctest.Wait()
			listener.AssertExpectations(t)
		})
	})
	t.Run("handle succeeds", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			conn := mocks.NewMockConn(t)
			conn.On("Read", mock.Anything).Return(0, io.EOF).Once()
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return(conn, nil).Once()
			listener.On("Accept").Return((net.Conn)(nil), os.ErrClosed).Once()

			assert.ErrorIs(t, os.ErrClosed, Serve(listener, func(c net.Conn) error {
				conn.Read(nil)
				return nil
			}))
			synctest.Wait()
			listener.AssertExpectations(t)
			conn.AssertExpectations(t)
		})
	})

	t.Run("handle fails", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			conn := mocks.NewMockConn(t)
			conn.On("Read", mock.Anything).Return(0, io.EOF).Once()
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return(conn, nil).Once()
			listener.On("Accept").Return((net.Conn)(nil), os.ErrClosed).Once()
			assert.ErrorIs(t, os.ErrClosed, Serve(listener, func(c net.Conn) error {
				c.Read(nil)
				return errors.New("some error")
			}))
			synctest.Wait()
			listener.AssertExpectations(t)
			conn.AssertExpectations(t)
		})
	})

	t.Run("handle fails then succeeds", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			conn1 := mocks.NewMockConn(t)
			conn1.On("Read", mock.Anything).Return(0, io.EOF).Once()
			conn2 := mocks.NewMockConn(t)
			conn2.On("Read", mock.Anything).Return(0, io.EOF).Once()
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return(conn1, nil).Once()
			listener.On("Accept").Return(conn2, nil).Once()
			listener.On("Accept").Return((net.Conn)(nil), os.ErrClosed).Once()
			assert.ErrorIs(t, os.ErrClosed, Serve(listener, func(c net.Conn) error {
				c.Read(nil)
				if c == conn1 {
					return errors.New("some error")
				}
				if c == conn2 {
					return nil
				}
				t.Error("unexpected connection received")
				return nil
			}))
			synctest.Wait()
			listener.AssertExpectations(t)
			conn1.AssertExpectations(t)
			conn2.AssertExpectations(t)
		})
	})

	t.Run("handle both connections takes one second", func(t *testing.T) {
		delay := time.Second
		var done atomic.Bool
		synctest.Test(t, func(t *testing.T) {
			// initialize the connections
			conn1 := mocks.NewMockConn(t)
			conn1.On("Read", mock.Anything).Run(func(args mock.Arguments) {
				time.Sleep(delay)
			}).Return(0, io.EOF).Once()
			conn2 := mocks.NewMockConn(t)
			conn2.On("Read", mock.Anything).Run(func(args mock.Arguments) {
				time.Sleep(time.Second)
			}).Return(0, io.EOF).Once()

			// set up the listener to accept both connections
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return(conn1, nil).Once()
			listener.On("Accept").Return(conn2, nil).Once()
			listener.On("Accept").Return((net.Conn)(nil), os.ErrClosed).Once()

			// run Serve
			go func() {
				ts := time.Now()
				assert.Error(t, Serve(listener, func(c net.Conn) error {
					c.Read(nil)
					return nil
				}))
				assert.LessOrEqual(t, time.Since(ts), delay)
				done.Store(true)
			}()

			time.Sleep(time.Hour)
			synctest.Wait()
			if !done.Load() {
				t.Error("waited for Serve to finish for 1 hour!!! Why is it not done yet?")
			}

			listener.AssertExpectations(t)
			conn1.AssertExpectations(t)
			conn2.AssertExpectations(t)
		})
	})

	t.Run("waits for in-flight handlers before returning", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			waitCh := make(chan struct{})
			conn := mocks.NewMockConn(t)
			listener := mocks.NewMockListener(t)
			listener.On("Accept").Return(conn, nil).Once()
			listener.On("Accept").Return((net.Conn)(nil), os.ErrClosed).Once()

			done := false
			go func() {
				assert.Error(t, Serve(listener, func(c net.Conn) error {
					<-waitCh
					return nil
				}))
				done = true
			}()
			synctest.Wait()
			assert.False(t, done)
			close(waitCh)
			synctest.Wait()
			assert.True(t, done)
			listener.AssertExpectations(t)
			conn.AssertExpectations(t)
		})
	})
}
