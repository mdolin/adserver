## Ad Server

The Ad Server is a simple server application that handles ad requests and selects the most applicable creative based on the AdUnit's format, width, and height.

### Main bits of the project
* Main File
* Database Package
* Handlers Package
* Models Package

### Structure of the project
```
.
├── README.md
├── ad.db
├── db
│   └── ad_db.go
├── go.mod
├── go.sum
├── handlers
│   └── ad_handlers.go
├── main.go
└── models
    └── entities.go
```

### Features
* Ad request handling: The server receives ad requests, retrieves the specified ad unit, filters through available creatives, and selects the most relevant one based on format, width, and height.

* Database interaction: The project uses an SQLite database to store and retrieve ad units and creatives.

* Cache refreshing: The server periodically updates its cache with the latest data from the database to ensure up-to-date information for ad selection.

* Applicable creative selection: The server selects the creative with the highest price that matches the requested ad unit's format, width, and height.

### Endpoints

- **/adrequest**:  
  Submit an ad request by sending a JSON payload in the request body. The payload should include the `ad_unit_id` and `user_id` fields. The server retrieves the specified ad unit, filters the available creatives, selects the most relevant one based on format, width, and height, and returns the selected creative's details as a JSON response.

- **/refresh**:  
  Trigger a cache refresh from the database. This endpoint updates the server's cache with the latest data for ad units and creatives. It is useful for updating the server's cache when changes are made to the database outside of normal ad requests.


### Prerequisites

- Go (version 1.16 or higher)
- SQLite (version 3 or higher)

### Installation

1. Clone the repository:

   ```shell
   git clone <repository-url>

2. Install the dependencies:

   ```shell
   go mod download

3. Create the SQLite database:

   ```shell
   sqlite3 ad.db < schema.sql

### Usage
1. Start the server:
   ```shell
   go run main.go

2. Make ad request using curl or any other HTTP client:
   ```shell
   curl -X POST -H "Content-Type: application/json" -d '{"ad_unit_id": "<ad-unit-id>", "user_id": "<user-id>"}' http://localhost:8080/adrequest

Replace <ad-unit-id> with the ID of the AdUnit you want to request an ad for, and <user-id> with the ID of the user making the request.

The server will respond with a JSON containing the selected creative's details:

```json
{
  "creative_id": "<creative-id>",
  "content": "<creative-content>",
  "price": "<creative-price>",
  "user_id": "<user-id>"
}
```

### Usage
The server is configured to run on `http://localhost:8080` by default. If you want to change the host or port, modify the main.go file.
