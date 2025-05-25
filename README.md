# PengTune

Manual PID tuning for controllers in Peng.
Offers plumbing to support machine tuning in the future.

## About

This app was built from the official Wails Svelte template.
It is a go backend with a Svelte frontend, and uses NATS.io as a message broker
between the backend and Peng.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
