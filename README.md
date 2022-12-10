# Toronto Water Exporter
Exports water consumption from Toronto Water to an InfluxDB.

## Configuration
The configuration file must be a valid YAML file. Its path can be passed into the application as an argument, else **config.yml** is assumed.

Example **config.yml** file:
```
influxDB:
  url: http://192.168.0.252:9086
  token: abc
  organization: home
  bucket: torontowater
torontoWater:
  accountNumber: "1234"
  lastName: "DOE"
  postalCode: "ABC"
  lastPaymentMethod: 1
sleepDuration: 720
lookDaysInPast: 1
```

| Name                           | Description                                                                 |
|--------------------------------|-----------------------------------------------------------------------------|
| influxDB.url                   | address of InfluxDB2 server                                                 |
| influxDB.token                 | auth token to access InfluxDB2 server                                       |
| influxDB.organization          | organization of InfluxDB2 server                                            |
| influxDB.bucket                | name of bucket                                                              |
| torontoWater.accountNumber     | used to log into Toronto Water                                              |
| torontoWater.lastName          | used to log into Toronto Water                                              |
| torontoWater.postalCode        | used to log into Toronto Water                                              |
| torontoWater.lastPaymentMethod | used to log into Toronto Water                                              |
| sleepDuration                  | sleep time between exports in minutes, zero means run only once             |
| lookDaysInPast                 | how many days of the past should be considered                              |

## Docker
The exporter was written with the intent of running it in docker. You can also run it directly if this is preferred.

### Build Image
Execute following statement, then either start via docker or docker compose.
```
docker build -t toronto-water-exporter .
```

### Docker
```
docker run -d --restart unless-stopped --name=toronto-water-exporter -v ./config.yml:/config.yml toronto-water-exporter
```

### Docker Compose
```
version: "3.4"
services:
  toronto-water-exporter:
    image: toronto-water-exporter
    container_name: toronto-water-exporter
    restart: unless-stopped
    volumes:
      - ./config.yml:/config.yml:ro
```