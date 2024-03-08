# Redis ... in go!

## Overview

This was a project to get myself more familiar with Go. It is a very interesting language (per tutorialspoint)

> Go is a general-purpose language designed with systems programming in mind. It was initially developed at Google in the year 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. It is strongly and statically typed, provides inbuilt support for garbage collection, and supports concurrent programming.

Redis is another technology that interests me. Redis is an in memory database that is often used as a caching layer (but it can be used as a full blown database). In its simplest form it is a key value store.

## Features

This redis implementation supports a handful of commands (more are on the way)

- ping (sends PONG back)
- echo <input> (sends back <input>)
- set <key> <value> (set a key value pair in the db)
- set <key> <value> px <millis> (set a key value pair in the db that will expire after <millis> milliseconds)
- get <key> (gets the key from the db)

## Future

I hope to implement more features on this project in the near future including:

- data replication
- streams
- persistence