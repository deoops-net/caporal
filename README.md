# Caporal 

a remote api server for schedule docker jobs.

this can work with [dotam](https://github.com/deoops-net/dotam) smoothly

## Usage

```bash
docker pull deoops-net/caporal:latest

// to turn on auth
// set env AUTH_USER and AUTH_PASS
// other wise just leave them empty
// and you can implement auth method at proxy level
docker run -d -p 8080:8080 --name caporal -v /var/run/docker.sock:/var/run/docker.sock deoops-net/caporal
```

## Setttings

caporal accepts environment variables for settings

### NOT_PRIVATE
if this is set as true, caporal will not do authentication for the image registry

### AUTH_USER AUTH_PASSWORD
this pair env variable is used for auth dotam client the docker.auth section of your Dotamfile

### REG_USER REG_PASSWORD
if this pair is set, caporal will use this for registry authorization 

NOTE: when NOT_PRIVATE is not set, the registry authorization order is.
look for REG_ pair first then AUTH_ pair
