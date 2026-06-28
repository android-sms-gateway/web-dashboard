<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache License][license-shield]][license-url]
[![Go][go-shield]][go-url]
[![Svelte][svelte-shield]][svelte-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/android-sms-gateway/web-dashboard">
    <img src="https://raw.githubusercontent.com/android-sms-gateway/brand/master/logo.png" alt="Logo" width="120" height="120">
  </a>

<h3 align="center">SMSGate Web Dashboard</h3>

  <p align="center">
    Web-based administrative dashboard for the SMSGate ecosystem — turn Android devices into SMS gateways.
    <br />
    <a href="https://github.com/android-sms-gateway/web-dashboard"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/web-dashboard/issues">Report Bug</a>
    ·
    <a href="https://github.com/android-sms-gateway/web-dashboard/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Features](#features)
  - [Built With](#built-with)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Installation](#installation)
    - [Docker (recommended)](#docker-recommended)
    - [GitHub Releases](#github-releases)
    - [Build from source](#build-from-source)
  - [Configuration](#configuration)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Available Commands](#available-commands)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->
## About The Project

The SMSGate Web Dashboard is a standalone Go HTTP server with an embedded Svelte 5 single-page application. It provides a graphical user interface and REST API for managing Android devices configured as SMS gateways through the [SMSGate](https://github.com/android-sms-gateway) platform.

The dashboard acts as a proxy to the [SMSGate 3rd Party API](https://api.sms-gate.app/3rdparty/v1), handling authentication, session management, and real-time event streaming via Server-Sent Events.

### Features

- **Dashboard** — aggregated statistics: devices online/active/total, messages sent/pending/failed
- **Message Management** — paginated message list with filtering (by state, device, date range); send new SMS; view delivery status and timeline
- **Device Management** — view registered devices with online/offline status, remove devices
- **Webhook Management** — create, list, and delete webhooks for SMS events (received, sent, delivered, failed, MMS, data SMS, system ping)
- **API Token Management** — generate JWT tokens with granular scope selection (15 permission levels); copy and revoke tokens
- **Device Settings** — configure message intervals, SIM selection mode, processing order, ping intervals, log lifetime, webhook signing key, retry policy, gateway cloud URL, and encryption passphrase
- **Real-time Notifications** — live SSE stream for incoming messages, state changes, and device status updates with toast notifications
- **OpenAPI Documentation** — Swagger UI at `/api/v1/docs`
- **Prometheus Metrics** — metrics endpoint for monitoring

### Built With

* [![Go][go-shield]][go-url]
* [![Fiber][fiber-shield]][fiber-url]
* [![Fx][fx-shield]][fx-url]
* [![Svelte][svelte-shield]][svelte-url]
* [![Vite][vite-shield]][vite-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ARCHITECTURE -->
## Architecture

```text
Browser (Svelte SPA)  ──►  Go Dashboard (:3000)  ──►  SMSGate 3rd Party API (:443)
                                    │
                          ┌─────────┼─────────┐
                          │         │         │
                     Cookie     SSE      Webhook
                     Sessions   Events   Callbacks
```

- **Authentication**: credentials are validated against the SMSGate API, then stored in cookie-backed server sessions
- **Real-time**: the dashboard registers webhooks on the SMSGate server for each user session; incoming webhook callbacks are translated to SSE events and pushed to the browser
- **Frontend**: Svelte 5 SPA compiled to static files, embedded in the Go binary via `embed.FS`

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Installation

#### Docker (recommended)

Multi-platform Docker images are published to GitHub Container Registry:

```sh
docker pull ghcr.io/android-sms-gateway/web-dashboard:latest
```

Run with:

```sh
docker run -p 3000:3000 \
  -e HTTP__LISTEN=0.0.0.0:3000 \
  -e GATEWAY__URL=https://api.sms-gate.app/3rdparty/v1 \
  -e GATEWAY__WEBHOOK_URL=https://your-public-url/api/webhooks/callback \
  ghcr.io/android-sms-gateway/web-dashboard:latest
```

Images support `linux/amd64` and `linux/arm64`, are based on Alpine Linux, and run as a non-root user.

#### GitHub Releases

Pre-built binaries for Linux, macOS, and Windows are available on the [GitHub Releases page](https://github.com/android-sms-gateway/web-dashboard/releases).

```sh
# Example — download and run the Linux amd64 binary
curl -LO https://github.com/android-sms-gateway/web-dashboard/releases/latest/download/web-dashboard_linux_amd64.tar.gz
tar xzf web-dashboard_linux_amd64.tar.gz
./web-dashboard
```

#### Build from source

1. Clone the repository.
   ```sh
   git clone https://github.com/android-sms-gateway/web-dashboard.git
   cd web-dashboard
   ```
2. Install Go dependencies.
   ```sh
   make deps
   ```
3. Build the project (generates OpenAPI docs and compiles frontend).
   ```sh
   make build
   ```

### Configuration

The application is configured via environment variables or an optional YAML file:

| Variable               | Default                                       | Description                              |
| ---------------------- | --------------------------------------------- | ---------------------------------------- |
| `HTTP__LISTEN`         | `127.0.0.1:3000`                              | HTTP server bind address                 |
| `GATEWAY__URL`         | `https://api.sms-gate.app/3rdparty/v1`        | SMSGate 3rd Party API endpoint           |
| `GATEWAY__WEBHOOK_URL` | `http://localhost:3000/api/webhooks/callback` | Public callback URL for webhook events   |
| `CONFIG_PATH`          | —                                             | Path to optional YAML configuration file |

Example:
```sh
export HTTP__LISTEN=0.0.0.0:8080
export GATEWAY__URL=https://api.sms-gate.app/3rdparty/v1
./web-dashboard
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- DEVELOPMENT -->
## Development

Start the development server with live reload:

```sh
make air
```

This runs the Go server via `air` (hot reload on `.go` changes) and automatically builds the Svelte frontend. The Vite dev server proxies `/api` requests to the Go server.

### Prerequisites

* Go 1.25+
  ```sh
  go version
  ```
* Node.js 20+
  ```sh
  node --version
  ```
* `golangci-lint` (optional, for linting)
  ```sh
  golangci-lint version
  ```
* `air` (optional, for live reload)
  ```sh
  go install github.com/air-verse/air@latest
  ```

### Available Commands

| Command         | Description                                                                |
| --------------- | -------------------------------------------------------------------------- |
| `make build`    | Full production build (gen + compile Go binary)                            |
| `make test`     | Run tests with race detection                                              |
| `make coverage` | Run tests and generate coverage report                                     |
| `make lint`     | Run golangci-lint                                                          |
| `make fmt`      | Format code with golangci-lint                                             |
| `make air`      | Start development server with live reload                                  |
| `make gen`      | Generate OpenAPI docs and compile frontend assets (prerequisite for build) |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- API ENDPOINTS -->
## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Route                            | Auth | Description                                   |
| ------ | -------------------------------- | ---- | --------------------------------------------- |
| POST   | `/auth/login`                    | No   | Authenticate with SMSGate credentials         |
| POST   | `/auth/logout`                   | No   | Destroy session                               |
| GET    | `/auth/me`                       | Yes  | Get authenticated username                    |
| GET    | `/stats`                         | Yes  | Aggregated dashboard statistics               |
| GET    | `/messages`                      | Yes  | List messages (paginated, filterable)         |
| POST   | `/messages`                      | Yes  | Send a new SMS                                |
| GET    | `/messages/:id`                  | Yes  | Get message details with delivery timeline    |
| GET    | `/devices`                       | Yes  | List registered devices                       |
| DELETE | `/devices/:id`                   | Yes  | Remove a device                               |
| GET    | `/webhooks`                      | Yes  | List webhooks                                 |
| POST   | `/webhooks`                      | Yes  | Create a webhook                              |
| DELETE | `/webhooks/:id`                  | Yes  | Delete a webhook                              |
| GET    | `/settings`                      | Yes  | Get device settings                           |
| PATCH  | `/settings`                      | Yes  | Update device settings                        |
| POST   | `/tokens`                        | Yes  | Generate an API token                         |
| DELETE | `/tokens/:jti`                   | Yes  | Revoke an API token                           |
| GET    | `/events`                        | Yes  | SSE real-time event stream                    |
| POST   | `/api/webhooks/callback/:userId` | No   | Webhook callback receiver (called by SMSGate) |

OpenAPI documentation is available at `/api/v1/docs` when enabled.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Maintainer: [@capcom6](https://github.com/capcom6)

Project Link: [https://github.com/android-sms-gateway/web-dashboard](https://github.com/android-sms-gateway/web-dashboard)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/web-dashboard.svg?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/web-dashboard/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/web-dashboard.svg?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/web-dashboard/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/web-dashboard.svg?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/web-dashboard/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/web-dashboard.svg?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/web-dashboard/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/web-dashboard.svg?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/web-dashboard/blob/master/LICENSE
[go-shield]: https://img.shields.io/badge/go-1.25%2B-00ADD8?style=for-the-badge&logo=go
[go-url]: https://go.dev/
[fiber-shield]: https://img.shields.io/badge/Fiber-v2-00b894?style=for-the-badge
[fiber-url]: https://github.com/gofiber/fiber
[fx-shield]: https://img.shields.io/badge/Uber%20Fx-DI-6f42c1?style=for-the-badge
[fx-url]: https://github.com/uber-go/fx
[svelte-shield]: https://img.shields.io/badge/Svelte-5-FF3E00?style=for-the-badge&logo=svelte
[svelte-url]: https://svelte.dev/
[vite-shield]: https://img.shields.io/badge/Vite-8-646CFF?style=for-the-badge&logo=vite
[vite-url]: https://vitejs.dev/
