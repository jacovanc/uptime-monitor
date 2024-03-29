# Setup
The sqlite3 dependency requires CGO_ENABLED=1, which means we need gcc installed on the machine to compile C code. For windows the easiest way is to install [MSYS2](https://www.msys2.org/)

## Config
Create a .env file (can copy the .example-env file and edit). Update all with appropriate settings - they are all required.

## Build
```make
make build
```

## Run
```make 
make run
```

## Test
```bash
go test ./...
```

# TODO List

This project involves creating a website uptime monitor with a dashboard, using Go for backend development and SQLite for database management. The monitor will track the status and latency of websites, and display historical data in chart form.

## Setup and Initial Development
- [x] Set up Go development environment and base structure.
- [x] Initialize SQLite database with SQLx 
    - [x] Database file should be automatically created if it doesn't exist
    - [x] Sort issues for using SQLite3 (requires gcc)

## Monitoring Logic
- [x] Implement a function to check website status and latency using Go's `http` package.
- [x] Write logic to store the status and latency results in SQLite database.

## Scheduling and Regular Checks
- [x] Implement a scheduler in Go to check websites at regular intervals.
- [x] Ensure the monitoring process runs continuously or as a background task.

## Alerts
- [x] Alert downtime or high latency to configured email address(es) using Mailgun
- [x] Setup mock implementation of the EmailSender implementation for testing purposes 
- [x] Batch alerts email sends that happen in a short space of time

## Improve concurrency usage, options:
- [x] Use different goroutines per website rather than per ping. Ensures that the pings for each website are spaced out correctly regardless of how long they take.

## Config
- [x] Add configuration file for defining the websites to monitor, the emails to notify, and perhaps trigger limits for things like latency/downtime length for firing an alert

## Testing
- [ ] Implement useful tests. 

## Web Dashboard Development
- [ ] Set up a web server using a Go web framework (e.g., Gin or Echo).
- [ ] Develop routes and handlers for the dashboard and website details.
- [ ] Implement frontend templates using HTML/CSS.

## Chart Integration and Data Visualization
- [ ] Integrate a JavaScript charting library (e.g., Chart.js) for historical data display.
- [ ] Implement frontend logic to fetch and display charts based on historical data.

## Frontend Design and User Interface
- [ ] Design a user-friendly interface for the dashboard.
- [ ] Ensure clear display of website statuses and easy navigation to historical data.

## Deployment
- [x] Deploy the application on a suitable platform (ideally something simple like Heroku - just ensure that any PaaS does not sleep after x time). Perhaps DigitalOcean App Platform?

## Further
 - [ ] Consider how to prevent internet outage on the uptime-monitor host from marking tracked websites as down. (Server downtime is easy as we just wont have any db rows for that time period, but if just the internet is down it will have rows marked as 'down'). 
