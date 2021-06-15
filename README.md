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

While probably not really built for this purpose that has never stopped me before. :)

Loki does meet all the requirements because I can write new series which can be subtracted from existing series to make corrections.

The idea is to log an entry in Loki every 5s, each entry then equates to 5s of time on that project.  

This allows using a `count_over_time` function to see how long time was spent on each project.

See below to see how we handle corrections with Loki as well as how the data can be visualized.

### Visualization

Mainly I want to visualize time spent on projects at various levels of aggregation.

* Day
* Week
* Month
* Quarter
* Year? 

### Dashboards

(Sorry I need to get these published)

The main query for using this data looks like this:

```
# All the counts are /12 because there is one entry every 5 seconds, therefore every entry represents 5s of time. 
# If we divide the count by 12 we can turn the result into minutes.
# e.g. 36 counts would be 36 5s blocks of time = 180s or 3minutes total
# 36/12 = 3 (another way to think about this is there are 12 5s intervals in a minute)

# If there were corrections, take the count of all entries where type="add" and subtract all the entries where type="sub"
(sum by (project) ((count_over_time({job="timefidget", type="add"} | logfmt | project != "" [$__interval])/12)) 
  - (sum by (project) (count_over_time({job="timefidget", type="sub"} | logfmt | project != "" [$__interval])/12))) 
# If there were no corrections just count all the entries where type="add"
or (sum by (project) (count_over_time({job="timefidget", type="add"} | logfmt | project != "" [$__interval])/12))
```

This query will output the correct time in minutes for each project based on the time window of the dashboard query, and because we use $__interval it should be correct no matter what time window you select.

## Running

**NOTE** When I first built this I hadn't built the libraries for sending data directly from Arduino to Loki so instead I had a small Go server in the middle.

Now that this [library](https://github.com/grafana/loki-arduino) exists the ESP32 can send directly to Loki without the need for anything in the middle.

Currently the `fidgserver` code all exists in this project but you don't need it anymore.

You can instead just checkout the `arduino/fidgobject` folder for the Arduino sketch.

## Making corrections

It is possible to make corrections although it's currently quite cumbersome...

**NOTE** I'm hopefully going to build this into a Grafana plugin to make this something you can do from the UI, but until then, I'm sorry....

First you will need to build the `markup` command line tool in `cmd/markup` with `make build-markup`

This tool works by sending events to Loki with timestamps matching the range you wish to make changes. It essentially writes entries 5s apart just like the actual device does.

Two types of corrections can be made:

1. subtracting incorrect events
1. adding events

### Subtracting events

```
cmd/markup/markup -from=2021-05-19T10:00:00-04:00 -to=2021-05-19T10:28:10-04:00 -project="1-1"
```

This will add "subtraction entries" for the `1-1` project for the provided time range, if you run this command as is it will just output what it will send.

To actually send the data

```
cmd/markup/markup -from=2021-05-19T10:00:00-04:00 -to=2021-05-19T10:28:10-04:00 -project="1-1" -write=true
```

### Adding events

We need to add a few more flags, `typeLabelVal` is confusing and I'm sorry but it will always be `-typeLabelVal=add` (what this does is add `type=add` label to the data)

`version` is a label which allows you to further correct for mistakes, I will explain this more below, but because the initial entry would have been created with `type=add` and so will these entries, we need to make a new stream so that Loki doesn't yell at us for out of order entries by trying to write data to an existing stream. That's what this label does.

```
cmd/markup/markup -from=2021-06-11T09:30:00-04:00 -to=2021-06-11T10:30:00-04:00 -project="1-1" -typeLabelVal=add -version=1
```

Same as before this command will only show you what it's going to do until you do:

```
cmd/markup/markup -from=2021-06-11T09:30:00-04:00 -to=2021-06-11T10:30:00-04:00 -project="1-1" -typeLabelVal=add -version=1 -write=true
```

### Version label

As briefly mentioned above, Loki will not allow you to add older entries to an existing stream. A stream is defined by the labels. 

When events are sent from the device they will have these labels:

`{job="timefidget",type="add"}`

When we make subtractions we send

`{job="timefidget",type="sub"}`

These are separate streams, so no problem, but when we go to add and we already have data we _will_ have a problem so we do this:

`{job="timefidget",type="add",version="1"}`

This makes a new stream but we don't really care about the version label it will be ignored in our dashboards.

This also works for subtractions if you made a mistake you can also do

```
cmd/markup/markup -from=2021-05-19T10:00:00-04:00 -to=2021-05-19T10:28:10-04:00 -project="1-1" -version=2 -write=true
```

You can increment the version label as much as you want to add more streams to fix mistakes, the dashboard queries will always take any values for `add` and subtract any values with the `sub` labels to show the resulting difference.
