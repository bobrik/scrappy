# Scrappy

Scrappy visualizes resource usage in your Mesos cluster. Currently it only
outputs text version, but a nice web interface can eventually be implemented
on top of it.

It'd be nice to see that information exposed in Mesos UI nicely too.

## Usage

Scrappy comes in a Docker container. Just point it to your Mesos master URL
to get resource usage breakdown.

```
docker run --rm bobrik/scrappy -u http://mesos-master-host:mesos-master-port
```

Sample report for you to know what to expect:

```
dev: 0.80 CPUs, 0.50GB RAM / 2.00 CPUs, 0.98GB RAM
  roles:
    - *: 0.80 CPUs, 0.50GB RAM / 2.00 CPUs, 0.98GB RAM
      tasks: 2
        - scrappy-example: 0.40 CPUs, 0.25GB RAM
        - scrappy-example: 0.40 CPUs, 0.25GB RAM

dev: 0.40 CPUs, 0.25GB RAM / 2.00 CPUs, 0.98GB RAM
  roles:
    - *: 0.40 CPUs, 0.25GB RAM / 2.00 CPUs, 0.98GB RAM
      tasks: 1
        - scrappy-example: 0.40 CPUs, 0.25GB RAM
```

The following options are available:

* `-u` Mesos URL to fetch data from.
* `-s` sort method, one of `host`, `cpu`, `mem`.
* `-r` reverse sort order.
* `-f` role name to filter on.

## Known issues

* Frameworks doing multiple roles at once are not supported. We assume that
each framework only owns one role and each task can only span one role.
This is fine most of the time.
