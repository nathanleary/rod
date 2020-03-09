package rod_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"path/filepath"
	"time"

	"github.com/ysmood/kit"
	"github.com/ysmood/rod"
	"github.com/ysmood/rod/lib/input"
)

func (s *S) TestClosePage() {
	page := s.browser.Page(srcFile("fixtures/click.html"))
	defer page.Close()
	page.Element("button")
}

func (s *S) TestPageContext() {
	p := s.page.Timeout(time.Minute).CancelTimeout().Cancel()
	_, err := p.CallE(`() => {}`, nil)
	s.Error(err)
}

func (s *S) TestRelease() {
	res, err := s.page.EvalE(false, "", `() => document`, nil)
	kit.E(err)
	s.page.Release(res.Get("result.objectId").String())
}

func (s *S) TestPageCall() {
	s.Greater(s.page.Call("DOM.getDocument", nil).Get("root.nodeId").Int(), int64(0))
}

func (s *S) TestWindow() {
	page := s.browser.Page(srcFile("fixtures/click.html"))
	defer page.Close()

	bounds := page.GetWindow()
	defer page.Window(
		bounds.Get("left").Int(),
		bounds.Get("top").Int(),
		bounds.Get("width").Int(),
		bounds.Get("height").Int(),
	)

	page.WindowMaximize()
	page.WindowNormal()
	page.WindowFullscreen()
	page.WindowNormal()
	page.WindowMinimize()
	page.WindowNormal()
	page.Window(0, 0, 1211, 611)
	s.EqualValues(1211, page.Eval(`() => window.innerWidth`).Int())
	s.EqualValues(611, page.Eval(`() => window.innerHeight`).Int())
}

func (s *S) TestSetViewport() {
	page := s.browser.Page(srcFile("fixtures/click.html"))
	defer page.Close()
	page.Viewport(317, 419, 0, false)
	res := page.Eval(`() => [window.innerWidth, window.innerHeight]`)
	s.EqualValues(317, res.Get("0").Int())
	s.EqualValues(419, res.Get("1").Int())

	page2 := s.browser.Page(srcFile("fixtures/click.html"))
	defer page2.Close()
	res = page2.Eval(`() => [window.innerWidth, window.innerHeight]`)
	s.NotEqual(int64(317), res.Get("0").Int())
}

func (s *S) TestPageLoadScript() {
	p := s.page.Navigate(srcFile("fixtures/click.html")).WaitLoad()

	res := p.LoadScript(srcFile("fixtures/load-script.js")).Eval(`() => ok()`)

	s.Equal("ok", res.String())
}

func (s *S) TestUntilPage() {
	page := s.page.Timeout(3 * time.Second).Navigate(srcFile("fixtures/open-page.html"))
	defer page.CancelTimeout()

	wait := page.WaitPage()

	page.Element("a").Click()

	newPage := wait()

	s.Equal("click me", newPage.Element("button").Text())

	wait()
}

func (s *S) TestPageWaitRequestIdle() {
	srv := kit.MustServer("127.0.0.1:0")
	defer func() { kit.E(srv.Listener.Close()) }()

	host := srv.Listener.Addr().String()

	srv.Engine.GET("/r", func(ctx kit.GinContext) { kit.Sleep(1) })
	srv.Engine.GET("/", ginHTML(`<html><button onclick="fetch('/r').then(r => r.text().then(t => t))">click</button></html>`))

	go func() { kit.Noop(srv.Do()) }()

	page := s.page.Navigate("http://" + host)

	page.Call("Network.enable", nil)

	wait := page.WaitRequestIdle()
	page.Element("button").Click()
	start := time.Now()
	wait()
	s.True(time.Since(start) > time.Second)
}

func (s *S) TestPageWaitIdle() {
	p := s.page.Navigate(srcFile("fixtures/click.html"))
	p.Element("button").Click()
	p.WaitIdle()

	s.True(p.Has("[a=ok]"))
}

func (s *S) TestPageWaitEvent() {
	wait := s.page.WaitEvent("Page.frameNavigated")
	s.page.Navigate(srcFile("fixtures/click.html"))
	wait()
}

func (s *S) TestAlert() {
	page := s.page.Navigate(srcFile("fixtures/alert.html"))

	wait := page.HandleDialog(true, "")

	go page.Element("button").Click()

	wait()
}

func (s *S) TestDownloadFile() {
	srv := kit.MustServer("127.0.0.1:0")
	defer func() { kit.E(srv.Listener.Close()) }()

	host := srv.Listener.Addr().String()
	content := "test content"

	srv.Engine.GET("/d", func(ctx kit.GinContext) {
		kit.E(ctx.Writer.Write([]byte(content)))
	})
	srv.Engine.GET("/", ginHTML(fmt.Sprintf(`<html><a href="//%s/d" download>click</a></html>`, host)))

	go func() { kit.Noop(srv.Do()) }()

	page := s.page.Navigate("http://" + host)

	wait := page.GetDownloadFile("*")

	page.Element("a").Click()

	_, data := wait()

	s.Equal(content, string(data))
}

func (s *S) TestMouse() {
	page := s.page.Navigate(srcFile("fixtures/click.html"))
	page.Element("button")
	mouse := page.Mouse

	mouse.Move(140, 160)
	mouse.Down("left")
	mouse.Up("left")

	s.True(page.Has("[a=ok]"))
}

func (s *S) TestMouseClick() {
	s.browser.Slowmotion(1)
	defer func() { s.browser.Slowmotion(0) }()

	page := s.page.Navigate(srcFile("fixtures/click.html"))
	page.Element("button")
	mouse := page.Mouse
	mouse.Move(140, 160)
	mouse.Click("left")
	s.True(page.Has("[a=ok]"))
}

func (s *S) TestDrag() {
	s.T().Skip("not able to use mouse event to simulate it for now")

	page := s.page.Navigate(srcFile("fixtures/drag.html"))
	mouse := page.Mouse

	mouse.Move(60, 30)
	mouse.Down("left")
	kit.E(mouse.MoveE(60, 80, 5))
	mouse.Up("left")

	page.Element(".dropzone:nth-child(2) #draggable")
}

func (s *S) TestPagePause() {
	go s.page.Pause()
	kit.Sleep(0.03)
	go s.page.Eval(`() => 10`)
	kit.Sleep(0.03)
	kit.E(s.page.CallE("Debugger.resume", nil))
}

func (s *S) TestPageScreenshot() {
	p := s.page.Navigate(srcFile("fixtures/click.html"))
	p.Element("button")
	data := p.Screenshot()
	img, err := png.Decode(bytes.NewBuffer(data))
	kit.E(err)
	s.Equal(800, img.Bounds().Dx())
	s.Equal(600, img.Bounds().Dy())
}

func (s *S) TestPageTraceDir() {
	p := *s.page.Navigate(srcFile("fixtures/click.html"))
	dir := filepath.FromSlash("tmp/trace-screenshots/" + kit.RandString(8))
	p.TraceDir(dir)
	p.Element("button").Click()
	pattern := filepath.Join(dir, "*")
	s.Len(kit.Walk(pattern).MustList(), 1)
}

func (s *S) TestPageInput() {
	p := s.page.Navigate(srcFile("fixtures/input.html"))

	el := p.Element("input")
	el.Focus()
	p.Keyboard.Press('A')
	p.Keyboard.InsertText(" Test")
	p.Keyboard.Press(input.Tab)

	s.Equal("A Test", el.Eval(`() => this.value`).String())
}

func (s *S) TestPageScroll() {
	p := s.page.Navigate(srcFile("fixtures/scroll.html")).WaitLoad()
	p.Mouse.Scroll(100, 200)
	kit.E(p.Mouse.ScrollE(200, 300, 5))
	p.Element("button").WaitStable()
	s.EqualValues(300, p.Eval("() => window.pageXOffset").Int())
	s.EqualValues(500, p.Eval("() => window.pageYOffset").Int())
}

func (s *S) TestPageOthers() {
	p := s.page.Navigate(srcFile("fixtures/input.html"))

	s.Equal("body", p.ElementByJS(`() => document.body`).Describe().Get("localName").String())
	s.Len(p.ElementsByJS(`() => document.querySelectorAll('input')`), 3)
	s.EqualValues(1, p.Eval(`() => 1`).Int())

	s.Panics(func() {
		rod.CancelPanic(errors.New("err"))
	})

	s.False(rod.IsError(io.EOF, rod.ErrElementNotFound))

	p.Mouse.Click("")
	p.Mouse.Down("left")
	defer p.Mouse.Up("left")
	p.Mouse.Down("right")
	defer p.Mouse.Up("right")
}

func (s *S) TestPageErrors() {
	p := s.page.Navigate(srcFile("fixtures/input.html"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.Context(ctx).NavigateE("")
	s.Error(err)

	err = p.Context(ctx).WindowE(nil)
	s.Error(err)

	_, err = p.Context(ctx).GetDownloadFileE("", "")
	s.Error(err)

	_, err = p.Context(ctx).ScreenshotE(nil)
	s.Error(err)

	err = p.Context(ctx).PauseE()
	s.Error(err)
}
