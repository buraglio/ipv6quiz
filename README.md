# IPv6 Quiz Application

A web-based quiz application that tests knowledge of IPv6 addressing and protocols, written in Go.

## Features

- 70+ IPv6-related multiple choice questions
- Random selection of 20 questions per quiz
- Web interface accessible via IPv6 loopback
- Scoring with percentage calculation
- Daemon mode for background operation
- Responsive HTML interface with clean styling

## Installation

1. Ensure you have Go installed (version 1.16+ recommended)
2. Clone this repository or download the source files
3. Build the application:

```
go build
```

## Basic usage 

```
./ipv6quiz
```

This starts the server in foreground mode, accessible at:
[http://[::1]:5000](http://[::1]:5000)

## Daemon mode

```
./ipv6quiz -d
```

The daemon will:
- Detach from the terminal
- Write its PID to ipv6quiz.pid
- Run until manually stopped

## Managing the Daemon

to kill the daemon

```
kill $(cat ipv6quiz.pid)
```

## Configuration
The application has the following default settings:

Bind address: [::1]:5000 (IPv6 loopback)
Number of questions per quiz: 20
Questions database: Hardcoded in questions.go
To change these, modify the source code before building.

## API Endpoints
GET / - Start a new quiz
POST /submit - Submit quiz answers
Building Questions
To add or modify questions, edit the questions slice in the source code. Each question follows this format:

```
{
    Text: &quot;Question text here&quot;,
    Options: []string{
        &quot;Option 1&quot;,
        &quot;Option 2&quot;,
        &quot;Option 3&quot;,
        &quot;Option 4&quot;,
    },
    Answer: indexOfCorrectAnswer, // 0-based
},
```

## Dependencies
Standard Go libraries only
No external dependencies required

## System Requirements
Operating System: Any supported by Go (Linux, macOS, Windows)
IPv6 stack enabled
Go 1.16+ (for building from source)

## Known Limitations
Currently only binds to IPv6 loopback, but why would you need anything else?
No persistent storage of questions or results
Basic authentication not implemented

## License

MIT License

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
