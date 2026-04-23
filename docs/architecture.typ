#import "@preview/cetz:0.3.4"

#set page(width: auto, height: auto, margin: 1cm)

#cetz.canvas({
  import cetz.draw: *

  // Colors
  let browser-color = rgb(135, 206, 250)   // light blue
  let server-color  = rgb(144, 238, 144)   // light green
  let bot-color     = rgb(255, 218, 185)   // peach
  let stun-color    = rgb(221, 160, 221)   // plum
  let arrow-color   = rgb(60, 60, 60)
  let text-color    = rgb(30, 30, 30)

  // --- Nodes ---

  // Browser (top-left)
  rect((-6, 2), (-2, 4), fill: browser-color, stroke: 1.5pt + text-color, radius: 0.3cm, name: "browser")
  content(("browser.center"), [
    #set text(font: "DejaVu Sans", size: 12pt, weight: "bold", fill: text-color)
    Browser Client
    #v(2pt)
    #set text(size: 9pt, weight: "regular")
    webrtc-web
    #linebreak()
    Next.js / React / TS
  ])

  // Signalling Server (top-right)
  rect((2, 2), (6, 4), fill: server-color, stroke: 1.5pt + text-color, radius: 0.3cm, name: "server")
  content(("server.center"), [
    #set text(font: "DejaVu Sans", size: 12pt, weight: "bold", fill: text-color)
    Signalling Server
    #v(2pt)
    #set text(size: 9pt, weight: "regular")
    webrtc-server
    #linebreak()
    Go / Gorilla WS
  ])

  // Bot Agents (bottom-left)
  rect((-6, -2), (-2, 0), fill: bot-color, stroke: 1.5pt + text-color, radius: 0.3cm, name: "bots")
  content(("bots.center"), [
    #set text(font: "DejaVu Sans", size: 12pt, weight: "bold", fill: text-color)
    Bot Agents
    #v(2pt)
    #set text(size: 9pt, weight: "regular")
    webrtc-agents
    #linebreak()
    Go / Pion WebRTC
  ])

  // STUN/TURN (bottom-right)
  rect((2, -2), (6, 0), fill: stun-color, stroke: 1.5pt + text-color, radius: 0.3cm, name: "stun")
  content(("stun.center"), [
    #set text(font: "DejaVu Sans", size: 12pt, weight: "bold", fill: text-color)
    STUN / TURN
    #v(2pt)
    #set text(size: 9pt, weight: "regular")
    coturn-deploy
    #linebreak()
    NAT Traversal
  ])

  // --- Arrows ---

  // 1. WebSocket signalling (browser <-> server)
  line((-2, 3), (2, 3), stroke: 1.5pt + arrow-color, mark: (end: ">", start: "<") )
  content((0, 3.6), [
    #set text(font: "DejaVu Sans", size: 9pt, fill: text-color)
    WebSocket Signalling
  ])
  content((0, 2.4), [
    #set text(font: "DejaVu Sans", size: 8pt, fill: rgb(100,100,100))
    SDP Offers / Answers / ICE Candidates
  ])

  // 2. WebSocket signalling (bots <-> server)
  line((-2, -1), (2, -1), stroke: 1.5pt + arrow-color, mark: (end: ">", start: "<") )
  content((0, -0.4), [
    #set text(font: "DejaVu Sans", size: 9pt, fill: text-color)
    WebSocket Signalling
  ])
  content((0, -1.6), [
    #set text(font: "DejaVu Sans", size: 8pt, fill: rgb(100,100,100))
    Register · SDP / ICE Relay
  ])

  // 3. WebRTC P2P (browser <-> bots) — dashed, crossing through the middle
  line((-4, 0.2), (-4, 1.8), stroke: (dash: "dashed", thickness: 1.5pt, paint: arrow-color), mark: (end: ">", start: "<") )
  content((-5.0, 1.0), [
    #set text(font: "DejaVu Sans", size: 8pt, fill: text-color)
    P2P Data
    #linebreak()
    Channels
  ])

  // 4. STUN/TURN assist (server <-> stun) — dotted
  line((4, 0.2), (4, 1.8), stroke: (dash: "dotted", thickness: 1.5pt, paint: arrow-color), mark: (end: ">", start: "<") )
  content((5.0, 1.0), [
    #set text(font: "DejaVu Sans", size: 8pt, fill: text-color)
    ICE
    #linebreak()
    Candidates
  ])

  // 5. Browser may also query STUN
  line((-3.5, -0.2), (2.5, -0.2), stroke: (dash: "dashed", thickness: 1pt, paint: rgb(150,150,150)), mark: (end: ">", start: "<") )
  content((-0.5, -0.6), [
    #set text(font: "DejaVu Sans", size: 7pt, fill: rgb(120,120,120))
    ICE / STUN queries (indirect)
  ])

  // --- Legend / Note ---
  rect((-6.5, -4.5), (6.5, -3), fill: rgb(250,250,250), stroke: 0.8pt + rgb(180,180,180), radius: 0.2cm, name: "legend")
  content(("legend.center"), [
    #set text(font: "DejaVu Sans", size: 9pt, fill: rgb(80,80,80))
    #strong[Note:] The signalling server only brokers connection setup. All chat messages,
    file transfers, pings, and audio tracks flow directly over WebRTC peer connections.
  ])
})
