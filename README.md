# open-says-me

This is a toy implementation of port knocking in Go.  Currently, it only supports iptables as the backend.

# Usage

Specify the port you wish to protect (`-port`) and several ports required to knock, in order, `-knock`.

```
Usage of open-says-me:
  -config string
    	config file (optional)
  -debug
    	log debug information
  -knock value
    	knock port (multiple supported)
  -port int
    	port to protect (default 9000)
  -pretty
    	pretty print logs
```

## Config File

This is make a little simpler with a config file and the `-config` flag.

```
port 9099
knock 8080
knock 8090
```

Protects `9099` by requiring knocks on `8080` then `8090`.

# Knocking

There is a client in `cmd/client` which can knock for you:

```
$ ./client 4000 5000
```

Alternatively you can use a tool like `netcat`:

```
$ echo -n "ok" | nc -u -c 127.0.0.1 4000
$ echo -n "ok" | nc -u -c 127.0.0.1 5000
```