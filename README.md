# Golang Uptime Slack bot

This application can monitor one or multiple websites and send notification to slack channel if the website is down.

## Configuration

### .env file
With this configuration you can only configure 1 domain
Create a file called ```.env``` into the same folder where your executable is

```
SLACK_BOT_TOKEN=<your token>
SLACK_CHANNEL_ID=<your channel id>

MONITOR_URL=https://yourdomain.com/health
MONITOR_TEXT="<html"

# Slow warning limit is set in miliseconds
SLOW_WARNING_LIMIT=3000

# Scan frequency in seconds
SCAN_FREQUENCY=30

# User agent, if not set it uses GolangUptimeBot/1.0
HTTP_USER_AGENT="TestUptimeBot/1.0"
```

### Yaml file
With this configuration you can only configure 1 or multiple domains.
Create a file called ```config.yaml``` into the same folder where your executable is

```
ScanFrequency: 0
Accounts:
- SlackBotToken: your token
  SlackChannelId: your channel id
  MonitorUrl: https://yourdomain1.com
  MonitorText: <html
  HttpUserAgent: UptimeBot/1.0
  SlowWarningLimit: 3000
- SlackBotToken: your token
  SlackChannelId: your channel id
  MonitorUrl: https://yourdomain2.com/health
  MonitorText: Welcome 
  HttpUserAgent: UptimeBot/1.0
  SlowWarningLimit: 6000
```

It is possible to set the same token and same channel, or same token different channel or different token and different channel in any combination
Note: The HttpUserAgent not required, if not set it will defult to `GolangUptimeBot/1.0` as of now

### Closing the applications
```
ctrl + c
```


## Make targets
```
make build
make run
make run-background
make run-test
```

The ```run-background``` will start the application (on linux and mac) in the backgound.

### Kill application running on background

You can kill your application if it is running in the background as follows:

```
ps -ax | grep upclient
```

Look from the process ID from the output:
```
77193 pts/0    Sl     0:00 ./upclient
77312 pts/0    S+     0:00 grep --color=auto upclient
```

Look for the one without the grep and copy the id
```
kill 77193
```

Application killed

## Coming soon
- Report with daily maximum, minimum and avarage page loads, number of outages once a day
- Database support to save report 
- Different scan frequency per user
- Avoid spamming of Slack if the error is too frequent and group them into one message per X time period.
- Sending out the daily report to slack as well
