# dummy

**dummy** is a command-line tool that turns Postman collections into functional mock APIs.
All mock endpoints and responses are taken directly from the example responses defined in the Postman collection.

## Warning :warning:
This project is still in WIP. :construction:
Expect rough edges and occasional bugs.

## Installation

### Prerequisites

-   Go 1.23+ recommended

### Build from source
```bash
git clone https://github.com/mathiasdonoso/dummy
cd dummy
go build -o dummy ./cmd/dummy
mv dummy ~/.local/bin/
```

## Usage

Run a mock server based on a Postman collection:
```bash
dummy run postman <path_to_postman_collection.json>
```

Example:
```bash
dummy run postman ./collections/auth.json
```

## What id does?
- Reads the collection’s example responses
- Creates mock endpoints based on each request in the collection
- Matches incoming requests by method + path + body
- Returns the exact status code, headers, and body from the Postman example

## Contributing
Feel free to open issues or submit pull requests to improve the tool. Contributions are welcome.
