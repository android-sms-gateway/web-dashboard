# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### New Features

#### Messages

- **Outgoing MMS support** — send multimedia messages via the API and web dashboard
  - New `mmsMessage` field in the send message API (use instead of `text`)
  - Exactly one of `text` or `mmsMessage` is required
  - New `deviceId` field to target a specific device
  - Attachment upload and preview in the compose page
  - Attachment display on the message detail page

### Documentation

- Updated README with screenshots for all dashboard pages
