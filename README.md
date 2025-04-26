# eventreminder

Simple Go application (script?) that reads in events from a CSV file and sends push notificatons if the event dates match the current date.

This was made for my amusement and entertainment only, to keep practising basic Go skills.

## Build

`go build -o eventreminder main.go`

## Usage

```bash
./eventreminder -file /path/to/events.csv
```

The application will accept a `-file` argument with a path to a CSV file, or if this is omitted it will try open a file named `events.csv` in the same folder as the application. The format of the csv file should be as follows:

```csv
day,month,year,subject,event,note
8,5,1958,Grandma Josephine and Grandpa Joe,Wedding Anniversary,don't call after 5pm
6,1,Peter Richardson,Birthday,
```

The eventreminder Go application requires the Pushover API token and Pushover recipient token in environment variables: `PUSHOVER_API_TOKEN` and `PUSHOVER_RECIPIENT_TOKEN`. 

I'm running this using a crontab entry, so I have a helper script `eventscript.sh` which sets these variables and then runs the Go application.
