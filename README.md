## Usage

```
  -baseCompose string
        The base docker compose file (default "docker-compose.default.yml")
  -composeDirName string
        Name of directory containing docker compose files (default ".compose")
  -debug
        Print debug messages
  -env string
        Environment that docker compose is running in. Will try to read the ENVIRONMENT env variable if not provided. (default "devel")
```

## Examples

```
> dc up -d
```
finds the nearest compose directory and executes the following
```
> docker-compose -f {compose_dir}/docker-compose.default.yml -f {compose_dir}/docker-compose.{env}.yml up -d
```

`env` can be overridden by an `ENVIRONMENT` environment variable
```
export ENVIRONMENT=staging
dc config --services
# will look for docker-compose.staging.yml
```
or by explicitly overriding
```
> dc --env=testing config
```
