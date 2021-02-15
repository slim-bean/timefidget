# TimeFidget

For people who like physical tools to interact with the digital world.

Inspired by [TimeFlip](https://timeflip.io/), in fact I ordered a TimeFlip2 but couldn't wait the 4 weeks for shipping so I decided to see if I could build a worse version myself.

## Goals

* Working prototype in a weekend. (Use parts on hand)
* Accessible for others. (Use easily obtainable parts)

## Requirements

* Move an object and have it start tracking time spent based on an upward facing surface.
* Be able to modify entries. (Mistakes happen)
* Visualize summaries in Grafana

## Non Goals

* 

## Design

### Storage

#### Prometheus

This was my first choice, however with no way to edit events it doesn't meet one of my requirements. In leu of editing data I consdired inserting time series which could be used to make corrections.

```
tracking_seconds_total{project=loki"} - tracking_seconds_corrected_total{project="loki"}
```

This works however the use is a little tricky
* How do you insert this corrected data? Do you insert one counter at some time near the correction with the value to subtract? What if the query time range doesn't include this correction?
* What if you want to insert it over an older time, Prometheus/Cortex don't really support this yet.

#### SQL

Probably my second choice and still a good fit, I have a lot of familiarty with SQL, Grafana already has built in support for creating queries, it's easy to edit entries, but maintaining relational databases over time is annoying. Additionally the queries can be cumbersome to write depending on how you store the data.

The biggest vote I have against SQL is I can't use SQLITE without CGO, and I really don't want to use CGO.

#### JSON

This is also not a bad fit, super simple, easy to edit files but this is a pretty bespoke option requiring me to write all the code to handle the files, rotating them and querying them. There is a JSON datasource over HTTP plugin available for Grafana, but it would require writing the query layer into these files. This isn't a bad option though, once the code is created at least the operations are very simple.

#### Loki

While probably not really built for this purpose that has never stopped me before. 

Loki 2.0 added features which would allow extracting numbers from log lines and aggregating in useful ways

```
sum(sum_over_time{project="loki",action="add"} | logfmt | unwrap duration [5m]) - sum(sum_over_time{project="loki",action="subtract"} | logfmt | unwrap duration [5m])
```

Instant queries would be best for this however, Loki doesn't currently split or shard instant queries so range queries will work better over longer time periods.

Using combinations of instant queries, range queries $__interval and reduce transforms in Grafana it should be possible to visualize the desired output.

### Visualization

Mainly I want to visualize time spent on projects at various levels of aggregation.

* Day
* Week
* Month
* Quarter
* Year? 

## Running

There are two ways to run the app:

1. Standalone, just logs the current active position every 5 seconds
2. Loki Embeded, a Loki server is started with the app and active position is sent to it every 5 seconds.

For the first mode, you would need to configure Promtail or another agent to ingest the logs and send them to a Loki instance.

Example:

```shell
./fidgserver [-port=8080]
```

For the second mode you need a Loki config file, an example is included:

Example:

```shell
./fidgserver -config.file=fidgserver-local-config.yaml
```

