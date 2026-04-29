// Package middlewares provides HTTP middleware components for the versitygw
// S3-compatible gateway API.
//
// # Authentication Middleware
//
// The auth middleware validates incoming requests against AWS Signature
// Version 4 (SigV4). It checks for the presence and basic format validity
// of the Authorization and X-Amz-Date headers.
//
// Usage:
//
//	validator := middlewares.NewSignatureValidator(accessKey, secretKey)
//	mux.Use(middlewares.AuthMiddleware(validator))
//
// On validation failure the middleware writes an S3-compatible XML error
// response and stops the request chain. Successful requests are forwarded
// to the next handler unchanged.
package middlewares
