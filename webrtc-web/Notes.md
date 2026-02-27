## Defined Environment Variables

This document describes the environment variables used in the WebRTC web application.

### `NEXT_PUBLIC_ICE_SERVERS`

- **Type:** String (comma-separated URLs)
- **Default:** `stun:stun.l.google.com:19302`
- **Used in:** `getICEServerURLs()` in `apis/ice.ts`
- **Description:** A comma-separated list of ICE server URLs used for WebRTC connection establishment. These servers are used by the main signaling server and the test server. If not set or empty, it defaults to Google's public STUN server.

**Example:**
```
NEXT_PUBLIC_ICE_SERVERS=stun:stun.l.google.com:19302,stun:stun1.l.google.com:19302
```

### `NEXT_PUBLIC_DN42_ICE_SERVERS`

- **Type:** String (comma-separated URLs)
- **Default:** Empty array `[]`
- **Used in:** `getDN42ICEServerURLs()` in `apis/ice.ts`
- **Description:** A comma-separated list of ICE server URLs specifically for DN42 and NeoNetwork connections. These servers are used by the DN42 signaling server when the application is accessed from a hostname ending in `.dn42`, `.neonetwork`, or `.neo`. If not set or empty, no ICE servers are configured for DN42 connections by default.

**Example:**
```
NEXT_PUBLIC_DN42_ICE_SERVERS=stun:stun.dn42.example.com:3478
```

## Notes

- Both environment variables are prefixed with `NEXT_PUBLIC_`, indicating they are exposed to the browser-side code in Next.js.
- The ICE servers are used for NAT traversal in WebRTC peer connections.
- URLs are trimmed and empty entries are filtered out when parsed.
- The signaling server selection logic in `apis/ws.ts` automatically prioritizes:
  - The "test" server when running on localhost
  - The "dn42" server when accessed from DN42 or NeoNetwork hostnames