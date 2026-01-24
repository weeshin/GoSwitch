# GoSwitch

GoSwitch is a robust and flexible ISO8583 engine written in Go.

## Overview

This library provides functionality to encode, decode, and manipulate ISO8583 messages. It is designed to be easy to use and extend for various payment system integrations.

## Features

- Parse and pack ISO8583 messages.
- Support for various field data types (numeric, alpha, binary, etc.).
- Customizable packager definitions.

## Installation

```bash
go get github.com/yourusername/GoSwitch
```

## Usage

```go
package main

import (
	"fmt"
	"GoSwitch/pkg/iso8583"
)

func main() {
	// Example usage
	fmt.Println("GoSwitch ISO8583 Engine")
}
```

## detailed documentation

Please refer to `pkg/iso8583` for more implementation details.
