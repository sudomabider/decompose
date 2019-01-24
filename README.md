```
> dc a -b -c
```
finds the nearest compose directory and executes the following
```
> docker-compose -f {compose_dir}/docker-compose.default.yml -f {compose_dir}/docker-compose.{env}.yml a -b -c
```

`env` defaults to `devel` and can be overridden by
```
> dc --env=testing a -b -c
```
