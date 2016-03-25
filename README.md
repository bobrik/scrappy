# Scrappy

Scrappy visualizes resource usage in your Mesos cluster. Currently it only
outputs text version, but a nice web interface can eventually be implemented
on top of it.

It'd be nice to see that information exposed in Mesos UI nicely too.

## Usage

Scrappy comes in a Docker container. Just point it to your Mesos master URL
to get resource usage breakdown. There are two flavors available.

### Slaves

`slaves` flavor reports stats on slave level:

```
docker run --rm bobrik/scrappy slaves -u http://mesos-master:port
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

### Roles

`roles` flavor reports stats on role level, useful if you have multiple of them:

```
docker run --rm bobrik/scrappy roles -u http://mesos-master:port
```

Sample report for just one default role for you to know what to expect:

```
- *: 0.30 CPUs, 0.12GB RAM / 2.00 CPUs, 0.98GB RAM
```

The following options are available:

* `-u` Mesos URL to fetch data from.
* `-f` role name to filter on.
