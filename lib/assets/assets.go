// generated by running "go generate" on project root

package assets

// Helper for rod
var Helper = `
(frameId) => { // eslint-disable-line no-unused-expressions
  const rod = {
    element (selector) {
      return (this.document || this).querySelector(selector)
    },

    elements (selector) {
      return (this.document || this).querySelectorAll(selector)
    },

    elementX (xpath) {
      return document.evaluate(
        xpath, (this.document || this), null, XPathResult.FIRST_ORDERED_NODE_TYPE
      ).singleNodeValue
    },

    elementsX (xpath) {
      const iter = document.evaluate(xpath, (this.document || this), null, XPathResult.ORDERED_NODE_ITERATOR_TYPE)
      const list = []
      let el
      while ((el = iter.iterateNext())) list.push(el)
      return list
    },

    elementMatches (selector, reg) {
      const r = new RegExp(reg)
      const filter = el => rod.text.call(el).match(r)
      const el = Array.from((this.document || this).querySelectorAll(selector)).find(filter)
      return el || null
    },

    parents (selector) {
      let p = this.parentElement
      const list = []
      while (p) {
        if (p.matches(selector)) {
          list.push(p)
        }
        p = p.parentElement
      }
      return list
    },

    async initMouseTracer (iconId, icon) {
      await rod.waitLoad()

      if (document.getElementById(iconId)) {
        return
      }

      const tmp = document.createElement('div')
      tmp.innerHTML = icon
      const svg = tmp.lastChild
      svg.id = iconId
      svg.style = 'position: absolute; z-index: 2147483647; width: 17px; pointer-events: none;'
      svg.removeAttribute('width')
      svg.removeAttribute('height')
      document.body.appendChild(svg)
      rod.updateMouseTracer(iconId, 0, 0)
    },

    updateMouseTracer (iconId, x, y) {
      const svg = document.getElementById(iconId)
      if (!svg) {
        return
      }
      svg.style.left = x - 2 + 'px'
      svg.style.top = y - 3 + 'px'
    },

    async overlay (id, left, top, width, height, msg) {
      await rod.waitLoad()

      const div = document.createElement('div')
      const msgDiv = document.createElement('div')
      div.id = id
      div.style = ` + "`" + `position: fixed; z-index:2147483647; border: 2px dashed red;
        border-radius: 3px; box-shadow: #5f3232 0 0 3px; pointer-events: none;
        box-sizing: border-box;
        left: ${left}px;
        top: ${top}px;
        height: ${height}px;
        width: ${width}px;` + "`" + `

      if (width * height === 0) {
        div.style.border = 'none'
      }

      msgDiv.style = ` + "`" + `position: absolute; color: #cc26d6; font-size: 12px; background: #ffffffeb;
        box-shadow: #333 0 0 3px; padding: 2px 5px; border-radius: 3px; white-space: nowrap;
        top: ${height}px;` + "`" + `

      msgDiv.innerHTML = msg

      div.appendChild(msgDiv)
      document.body.appendChild(div)

      if (window.innerHeight < msgDiv.offsetHeight + top + height) {
        msgDiv.style.top = -msgDiv.offsetHeight - 2 + 'px'
      }

      if (window.innerWidth < msgDiv.offsetWidth + left) {
        msgDiv.style.left = window.innerWidth - msgDiv.offsetWidth - left + 'px'
      }
    },

    async elementOverlay (id, msg) {
      const interval = 100

      let pre = rod.box.call(this)
      await rod.overlay(id, pre.left, pre.top, pre.width, pre.height, msg)

      const update = () => {
        const overlay = document.getElementById(id)
        if (overlay === null) return

        const box = rod.box.call(this)
        if (pre.left === box.left && pre.top === box.top && pre.width === box.width && pre.height === box.height) {
          setTimeout(update, interval)
          return
        }

        overlay.style.left = box.left + 'px'
        overlay.style.top = box.top + 'px'
        overlay.style.width = box.width + 'px'
        overlay.style.height = box.height + 'px'
        pre = box

        setTimeout(update, interval)
      }

      setTimeout(update, interval)
    },

    removeOverlay (id) {
      const el = document.getElementById(id)
      el && el.remove()
    },

    waitIdle (timeout) {
      return new Promise((resolve) => {
        window.requestIdleCallback(resolve, { timeout })
      })
    },

    waitLoad () {
      return new Promise((resolve) => {
        if (document.readyState === 'complete') return resolve()
        window.addEventListener('load', resolve)
      })
    },

    async scrollIntoViewIfNeeded () {
      if (!this.isConnected) { throw new Error('Node is detached from document') }
      if (this.nodeType !== Node.ELEMENT_NODE) { throw new Error('Node is not of type HTMLElement') }

      const visibleRatio = await new Promise(resolve => {
        const observer = new IntersectionObserver(entries => {
          resolve(entries[0].intersectionRatio)
          observer.disconnect()
        })
        observer.observe(this)
      })
      if (visibleRatio !== 1.0) { this.scrollIntoView({ block: 'center', inline: 'center', behavior: 'instant' }) }
    },

    inputEvent () {
      this.dispatchEvent(new Event('input', { bubbles: true }))
      this.dispatchEvent(new Event('change', { bubbles: true }))
    },

    selectText (pattern) {
      const m = this.value.match(new RegExp(pattern))
      if (m) {
        this.setSelectionRange(m.index, m.index + m[0].length)
      }
    },

    selectAllText () {
      this.select()
    },

    select (selectors) {
      selectors.forEach(s => {
        Array.from(this.options).find(el => {
          try {
            if (el.innerText.includes(s) || el.matches(s)) {
              el.selected = true
              return true
            }
          } catch (e) { }
        })
      })
      this.dispatchEvent(new Event('input', { bubbles: true }))
      this.dispatchEvent(new Event('change', { bubbles: true }))
    },

    visible () {
      const box = this.getBoundingClientRect()
      const style = window.getComputedStyle(this)
      return style.display !== 'none' &&
        style.visibility !== 'hidden' &&
        !!(box.top || box.bottom || box.width || box.height)
    },

    invisible () {
      return !rod.visible.apply(this)
    },

    box () {
      const box = this.getBoundingClientRect().toJSON()
      if (this.tagName === 'IFRAME') {
        const style = window.getComputedStyle(this)
        box.left += parseInt(style.paddingLeft) + parseInt(style.borderLeftWidth)
        box.top += parseInt(style.paddingTop) + parseInt(style.borderTopWidth)
      }
      return box
    },

    text () {
      switch (this.tagName) {
        case 'INPUT':
        case 'TEXTAREA':
          return this.value
        default:
          return this.innerText
      }
    },

    resource () {
      return new Promise((resolve, reject) => {
        if (this.complete) {
          return resolve(this.currentSrc)
        }
        this.addEventListener('load', () => resolve(this.currentSrc))
        this.addEventListener('error', (e) => reject(e))
      })
    },

    stripHTML (html) {
      const div = document.createElement('div')
      div.innerHTML = html
      return div.innerText
    },

    addScriptTag (id, url, content) {
      if (document.getElementById(id)) return

      return new Promise((resolve, reject) => {
        var s = document.createElement('script')

        if (url) {
          s.src = url
          s.onload = resolve
        } else {
          s.type = 'text/javascript'
          s.text = content
          resolve()
        }

        s.id = id
        s.onerror = reject
        document.head.appendChild(s)
      })
    },

    addStyleTag (id, url, content) {
      if (document.getElementById(id)) return

      return new Promise((resolve, reject) => {
        var el

        if (url) {
          el = document.createElement('link')
          el.rel = 'stylesheet'
          el.href = url
        } else {
          el = document.createElement('style')
          el.type = 'text/css'
          el.appendChild(document.createTextNode(content))
          resolve()
        }

        el.id = id
        el.onload = resolve
        el.onerror = reject
        document.head.appendChild(el)
      })
    }
  }

  window['rod' + frameId] = rod

  if (!window.rod) window.rod = rod

  return window
}
`

// MousePointer for rod
var MousePointer = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="277px" height="401px" viewBox="0 0 277 401" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <!-- Generator: Sketch 52.6 (67491) - http://www.bohemiancoding.com/sketch -->
    <title>mouse-pointer</title>
    <desc>Created with Sketch.</desc>
    <defs>
        <polygon id="path-1" points="0 0 0 299 60 241 103 341 170 313 130 218 217 216"></polygon>
        <filter x="-24.2%" y="-11.0%" width="148.4%" height="130.8%" filterUnits="objectBoundingBox" id="filter-2">
            <feOffset dx="0" dy="15" in="SourceAlpha" result="shadowOffsetOuter1"></feOffset>
            <feGaussianBlur stdDeviation="15" in="shadowOffsetOuter1" result="shadowBlurOuter1"></feGaussianBlur>
            <feColorMatrix values="0 0 0 0 0.138818027   0 0 0 0 0.138818027   0 0 0 0 0.138818027  0 0 0 0.502660779 0" type="matrix" in="shadowBlurOuter1"></feColorMatrix>
        </filter>
    </defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="mouse-pointer" transform="translate(30.000000, 15.000000)">
            <g id="outside">
                <use fill="black" fill-opacity="1" filter="url(#filter-2)" xlink:href="#path-1"></use>
                <use fill="#FFFFFF" fill-rule="evenodd" xlink:href="#path-1"></use>
            </g>
            <polygon id="inside" fill="#000000" points="18 44 18 255 66 207 110 313 145 299 102 198 171 197"></polygon>
        </g>
    </g>
</svg>`

// Monitor for rod
var Monitor = `<html>
<head>
    <title>Rod Monitor - Pages</title>
    <style>
        body {
            margin: 0;
            background: #2d2c2f;
            color: white;
            padding: 20px;
            font-family: sans-serif;
        }
        a {
            color: white;
            padding: 1em;
            margin: 0.5em 0;
            font-size: 1em;
            text-decoration: none;
            display: block;
            border-radius: 0.3em;
            border: 1px solid transparent;
            background: #212225;
        }
        a:visited {
            color: #c3c3c3;
        }
        a:hover {
            background: #25272d;
            border-color: #8d8d96; 
        }
    </style>
</head>
<body>
    <h3>Choose a Page to Monitor</h3>

    {{range .list}}
        <a href='/page/{{.TargetID}}' title="{{.URL}}">{{.Title}}</a>
    {{end}}
</body>
</html>`

// MonitorPage for rod
var MonitorPage = `<html>
<head>
    <title>Rod Monitor - {{.id}}</title>
    <style>
        body {
            margin: 0;
            background: #2d2c2f;
            color: #ffffff;
        }
        .navbar {
            font-family: sans-serif;
            border-bottom: 1px solid #1413158c;
            display: flex;
            flex-direction: row;
        }
        .error {
            color: #ff3f3f;
            background: #3e1f1f;
            border-bottom: 1px solid #1413158c;
            display: none;
            padding: 10px;
            margin: 0;
        }
        input {
            background: transparent;
            color: white;
            border: none;
            border: 1px solid #4f475a;
            border-radius: 3px;
            padding: 5px;
            margin: 5px;
        }
        .title {
            flex: 2;
        }
        .url {
            flex: 5;
        }
        .rate {
            flex: 1;
        }
    </style>
</head>
<body>
    <div class="navbar">
        <input type="text" class="title" title="title of the remote page" readonly>
        <input type="text" class="url" title="url of the remote page" readonly>
        <input type="number" class="rate" value="0.5" min="0" step="0.1" title="refresh rate (second)">
    </div>
    <pre class="error"></pre>
    <img class="screen">
</body>
<script>
    let elImg = document.querySelector('.screen')
    let elTitle = document.querySelector('.title')
    let elUrl = document.querySelector('.url')
    let elRate = document.querySelector('.rate')
    let elErr = document.querySelector('.error')

    async function update() {
        let res = await fetch('/api/page/{{.id}}')
        let info = await res.json()
        elTitle.value = info.title
        elUrl.value = info.url 

        await new Promise((resolve, reject) => {
            let now = new Date()
            elImg.src = '/screenshot/{{.id}}?t=' + now.getTime()
            elImg.style.maxWidth = innerWidth + 'px'
            elImg.onload = resolve
            elImg.onerror = () => reject('error loading screenshots')
        })
    }

    async function mainLoop() {
        try {
            await update()
            elErr.attributeStyleMap.delete("display")
        } catch (err) {
            elErr.style.display = "block"
            elErr.textContent = err + ""
        }

        setTimeout(mainLoop, parseFloat(elRate.value) * 1000)
    }

    mainLoop()
</script>
</html>`
