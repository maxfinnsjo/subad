# SUBAD — Agent Guide

## Project Overview

Subad is a lightweight, sandboxed virtual room platform. Users get a username and a starting room. Inside rooms they can place interactive items: decorations, posters, links, mirrors, windows, calendars, GIFs, videos, images, music play-buttons/backgrounds, and portals to other rooms.

The core philosophy: **simple, creative, unorthodox, secure, lightweight**.

## Tech Stack Decision

Keep it as simple as humanly possible:
- **Client**: Local HTML/JS app (single-file or minimal SPA)
- **Persistence**: Local storage + optional sync (CRDT or similar)
- **Crypto**: Optional WebCrypto for identity/rooms (no backend required)
- **Networking**: WebRTC / libp2p / lightweight relay for P2P room discovery
- **No database server, no cloud auth** — sandboxed, isolated by design

## Design Principles

1. **Sandboxed by default** — users are isolated until they choose to connect
2. **Organic social** — meet, share, create, subscribe, follow, display
3. **XP system** — rewards for active participation, room decoration, interaction
4. **Permission system** — room owners set who can enter, edit, or view
5. **Portals** — rooms connected through portals, creating a graph of spaces

## File Structure

```
subad/
├── AGENTS.md          # This file — dev guide
├── README.md          # Project intro
├── index.html         # Entry point (or src/ if it grows)
├── src/
│   ├── app.js         # Main app logic
│   ├── room.js        # Room management
│   ├── user.js        # User/identity module
│   ├── items/         # Room item types
│   │   ├── poster.js
│   │   ├── portal.js
│   │   ├── media.js
│   │   └── ...
│   ├── xp.js          # XP / reward system
│   ├── permissions.js # Permission model
│   └── network.js     # P2P / sync layer
├── assets/
│   └── ...
└── data/              # .gitkeep (persistence files live here)
```

## Development Workflow

- First pass: fully local (single-user, no network)
- Second pass: P2P sync between rooms
- Third pass: XP + permissions + portals graph

## Build & Run

No build step needed for v1. Open `index.html` in a browser or serve with:

```bash
python3 -m http.server 8080
```
