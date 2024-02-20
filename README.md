# Website Uptime Monitor Project TODO List

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
- [ ] Setup mock implementation of the EmailSender implementation for testing purposes

## Improve concurrency usage, options:
- [ ] Use a semaphore (challen) "ping limit" to ensure there is never more goroutines active than the number of websites (preventing overlap of the same website)?
- [ ] Add a timeout to ping requests lower than the sleep interval to prevent overlapping calls?
- [ ] Use different goroutines per website rather than per ping? Ensures that the pings for each website are spaced out correctly regardless of how long they take?
*I think preventing the same website from ever having overlapping goroutines is the right improvement. This prevents race conditions as overlapping goroutines only matters if they are accessing the same key on statusHistory (the website is the key). This also means that we never fire another ping for a website while still waiting on the result*

## Config
- [x] Add configuration file for defining the websites to monitor, the emails to notify, and perhaps trigger limits for things like latency/downtime length for firing an alert

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

## Testing
- [ ] Implement useful tests. 

## Deployment
- [ ] Deploy the application on a suitable platform (ideally something simple like Heroku - just ensure that any PaaS does not sleep after x time). Perhaps DigitalOcean App Platform?

## Further
 - [ ] Consider how to prevent internet outage on the uptime-monitor host from marking tracked websites as down. (Server downtime is easy as we just wont have any db rows for that time period, but if just the internet is down it will have rows marked as 'down'). 

# Setup
The sqlite3 dependency requires CGO_ENABLED=1, which means we need gcc installed on the machine to compile C code.

## Build
```make
make build
```

## Run
```make 
make run
```