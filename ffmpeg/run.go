// Copyright ? 2018 Wolfy-J <wolfy.jd@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ffmpeg

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type C struct {
	Name   string
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func New(name string, args []string) C {
	var c C
	c.Name = name
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.cmd = exec.CommandContext(c.ctx, "ffmpeg", args...)
	c.cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c.cmd.Stdout = os.Stdout
	return c
}

// Run runs ffmpeg with given set of arguments, optional callback will be used to report progress (current duration,
// total duration). Callback total duration can be 0 if unable to automatically detect.
func (c *C) Run() error {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//	cmd.Stdout = os.Stdout
	//	cmd.Start()

	//	time.Sleep(10 * time.Second)
	//	fmt.Println("退出程序中...", cmd.Process.Pid)
	//	cancel()

	//	return cmd.Wait()
	return c.cmd.Start()
}
func (c *C) RunThenClose(ch chan int) error {
	err := c.cmd.Start()
	c.cancel()
	//tell go i take it
	return err
}

func (c *C) Close() {
	c.cancel()
}

//run and wait for back
func RunAndClose(args []string, callback func(c, t time.Duration)) error {
	cmd := exec.Command("ffmpeg", args...)

	if callback == nil {
		var cmdErr bytes.Buffer
		cmd.Stderr = &cmdErr

		if err := cmd.Run(); err != nil {
			return extractError(err, cmdErr.String())
		}

		return nil
	}

	// ffmpeg stdout is stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	w, err := newWatcher(stderr, callback)
	if err != nil {
		return err
	}

	defer w.Close()
	defer cmd.Process.Wait()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
