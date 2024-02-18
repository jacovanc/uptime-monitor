# Website Uptime Monitor Project TODO List

This project involves creating a website uptime monitor with a dashboard, using Go for backend development and SQLite for database management. The monitor will track the status and latency of websites, and display historical data in chart form.

## Setup and Initial Development
- [x] Set up Go development environment and base structure.
- [ ] Initialize SQLite database with SQLx 
    - [ ] Database file should be automatically created if it doesn't exist
    - [ ] Sort issues for using SQLite3 (requires gcc)

#***REMOVED***ing Logic
- [x] Implement a function to check website status and latency using Go's `http` package.
- [x] Write logic to store the status and latency results in SQLite database.

## Scheduling and Regular Checks
- [x] Implement a scheduler in Go to check websites at regular intervals.
- [x] Ensure the monitoring process runs continuously or as a background task.

#***REMOVED***
- [ ] Alert downtime or high latency to configured email address(es) using Mailgun

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