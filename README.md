# Gonce
A HTTP API performance testing tool written in GoLang
* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)

## Description

A performance testing tool written in GoLang. Supports performance testing of HTTP APIs.

## Installation

The go application can be installed using:
```go get github.com/arham-jain/gonce```

Alternatively, the binary of the above can be downloaded [here](https://drive.google.com/open?id=1v_uqKNhKkB3YinCHgvTpSEnuyhG6Tl0j).

## Usage

Writing the configuration file. A sample config file:
```
{
    "threads" : 10,
    "rampUpTimeInSeconds":,
    "executionTimeInSeconds":5,
    "httpRequest": {
        "method":"POST",
        "urlWithEndpoint":"https://sample.com/v1/api",
        "payload": {
            {{sample_payload}}
        },
        "headers": {
            "Content-Type": "application/json"
        },
        "queryParams":{
            "param_key":"param_value
        },
        "successStatusCodes": {
            "201":true
        }
    }
}
```

Running the performance test:
```go run main.go {{path_to_file}}```

## Contribution

To contribute,

1. Create an issue on github
2. Create a pull request
3. Submit for review