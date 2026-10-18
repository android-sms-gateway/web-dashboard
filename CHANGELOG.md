# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Improvements

- **Docker healthcheck** — the published image now probes `/health/live` every 15 seconds after a 20-second start period, so Docker and orchestrators can detect an unhealthy container

## [0.6.3] - 2026-10-27

### Maintenance

- Updated `source-map-js` from `1.2.1` to `1.2.2`

## [0.6.2] - 2026-10-26

### Maintenance

- Upgraded the Go toolchain from `1.25.7` to `1.26.0`
- Updated `devalue` from `5.8.1` to `5.9.4`

## [0.6.1] - 2026-10-20

### New Features

#### Devices

- **End-to-end encryption status** — the devices table now shows whether a device uses end-to-end encryption, with the key rotation version when available
  - New optional `publicKey` and `keyVersion` fields in the devices API response; both are omitted for devices without end-to-end encryption configured

## [0.6.0] - 2026-10-14

### New Features

#### Messages

- **Outgoing MMS support** — send multimedia messages via the API and web dashboard
  - New `mmsMessage` field in the send message API (use instead of `text`)
  - Exactly one of `text` or `mmsMessage` is required
  - New `deviceId` field to target a specific device
  - Attachment upload and preview in the compose page
  - Attachment display on the message detail page

## [0.5.2] - 2026-09-26

### New Features

#### Messages

- **SMS parts counter** — the compose editor now shows a live character count, concatenated SMS part count, and encoding (GSM-7 or UCS-2) as you type

### Documentation

- Added a changelog, refreshed the README with screenshots for all dashboard pages, and updated `.env.example`

## [0.5.1] - 2026-08-22

### Bug Fixes

- Fixed messages with an empty `createdAt` timestamp breaking message views — they now display correctly

## [0.5.0] - 2026-08-20

### New Features

#### Webhooks

- **Batch webhook events** — the webhooks page now handles batch delivery events, displaying grouped message payloads instead of treating them as unknown events

## [0.4.1] - 2026-08-16

### Maintenance

- Updated dependencies (postcss security bump and latest module upgrades)

## [0.4.0] - 2026-08-07

### New Features

#### Dashboard

- **Trend charts** — 7/14/30-day message trend charts on the dashboard, backed by a pluggable cache (`memory://` or `redis://` via `CACHE__URL`)

## [0.3.0] - 2026-07-30

### New Features

#### Web Interface

- **Dark theme** — switch the dashboard between light and dark appearance

## [0.2.0] - 2026-07-20

### New Features

#### Web Interface

- **New single-page frontend** — ground-up rebuild of the dashboard UI
  - Login and authentication shell
  - Dashboard with real statistics and live SSE updates
  - Messages feature
  - Remaining screens migrated from the old UI

#### Observability

- **Grafana dashboard and Prometheus alerts** — SRE dashboard and alert rules shipped under `deployments/` for monitoring the dashboard service

## [0.1.2] - 2026-07-09

### Bug Fixes

- Fixed routing for URL-encoded path segments — percent-encoded paths now match routes correctly (framework update affecting the webhooks page)

## [0.1.1] - 2026-07-05

### Bug Fixes

- Fixed the live event stream dropping its connection unexpectedly — SSE updates now stay connected

## [0.1.0] - 2026-06-28

### New Features

#### Webhooks, Tokens & Settings

- **Webhooks** — create and manage webhook subscriptions for SMS events from the web UI
- **API tokens** — generate and revoke API tokens
- **Device settings** — configure device settings from the web UI

#### Real-time Updates

- **Live event stream** — the web UI processes Server-Sent Events for instant message, state, and device updates

#### Messages

- **Message pagination** — browse large message histories page by page

### Improvements

- Dashboard statistics now show current values
- Clearer error messages across the web UI
- Improved device selection and display
- General UI polish

## [0.0.1] - 2026-05-29

### New Features

#### API

- **BFF API** — backend-for-frontend REST API for managing devices and messages

#### Web Interface

- **Initial web dashboard** — first web UI with device and message views
