# Caporal 

a remote api server for schedule docker jobs.

this can work with [dotam](https://github.com/deoops-net/dotam) smoothly

## Usage

```bash
docker pull deoops-net/caporal

// to turn on auth
// set env AUTH_USER and AUTH_PASS
// other wise just leave them empty
// and you can implement auth method at proxy level
docker run -d -p 8080:8080 --name caporal -e AUTH_USER=tom -e AUTH_PASS=foo deoops-net/caporal
```