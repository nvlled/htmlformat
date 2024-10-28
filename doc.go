// Package htmlformat provides function for cleanly formatting html code,
// for purposes of making the html code more human-readable.
//
// Cleanly formatted (AKA pretty printed) means:
// - any excess or unnecessary whitespaces are removed
// - child nodes are indented and aligned relative to parent node
//
// This package does not handle long lines, no line
// wrapping is done to break lone lines.
//
// Conserving whitespaces takes priority over formatting and readability.
// This means this package avoids adding whitespaces where it might
// alter rendered HTML.
package htmlformat
