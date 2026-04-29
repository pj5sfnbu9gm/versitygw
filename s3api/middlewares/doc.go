// Package middlewares provides HTTP middleware components for the versitygw
// S3-compatible API gateway.
//
// Available middleware:
//
//   - AuthMiddleware: validates AWS Signature Version 4 request signatures.
//   - RequestIDMiddleware: attaches a unique request ID to each request context
//     and response header (X-Request-Id).
//   - LoggingMiddleware: structured access logging with latency, status code,
//     and request metadata.
//   - RateLimitMiddleware: per-IP sliding-window rate limiting that returns
//     HTTP 429 Too Many Requests when the configured threshold is exceeded.
//
// Middleware is designed to be composed with any standard net/http handler or
// router and does not impose ordering constraints beyond the documented
// dependency of LoggingMiddleware on RequestIDMiddleware for request-ID
// propagation in log entries.
package middlewares
