# dbtx

## Why

When we develop an app, we always depend on some databases. In the using of databases, we may use transactions to do something. Normally, we may transmit a variable named `tx` as transaction to each methods to do something. But, I think it's not elegant enough. So, `dbtx` is used to transmit transaction variables.

## How

Fundamentally, `dbtx` package `tx` to `context.Context`. We always transmit `context.Context` to each methods. So it can be more useful.
