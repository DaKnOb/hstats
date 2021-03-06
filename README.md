# hstats
A simple Web Server Log Analyzer

## What's that?
`hstats` is a tool written in Go to parse HTTP Request Logs from one or more
web servers. After the parsing is complete, the tool will output statistics
about the requests served by the server in either a human readable format or a
parseable one, for further use in scripts.

The initial version of the tool was created in 5-10 minutes in order to
help generate statistics really quick, without using any complex software.
It accepts logs that match some certain criteria from `stdin` and then
serially parses them to generate a report in the end. The current log
"parser" is a hack that has been thrown together and will be updated
soon (hopefully).

Future versions of this tool should be able to find the necessary information
in any given log line and possibly be NGiNX variable aware. Until then,
all Pull Requests and Issues are welcome.

## Can I use that?
If you see this text right here, in order to use `hstats` with the current
parser, your web server log files need to adhere to some "standards" so they
can be parsed properly. Luckily, if you haven't performed many changes, the
software will work "as-is". More specifically, your log files need to:

* They need to start with the IP Address (IPv4 or IPv6) of the remote client
* They need to contain the HTTP Protocol Version (HTTP/2, HTTP/1.1, HTTP/1)
  for the HTTP Protocol Version Feature to work.
* The HTTP Response Code must be stored as an integer, between two whitespaces
  (` `) and must be the element immediately following the second occurence of
  `"` in your log file line.

By default, the NGiNX and Apache log formats satisfy the above requirements.
Here are some excerpts of log lines in the above two web servers:

**NGiNX**:

```
2001:db8::5 - - [13/Dec/1989:13:37:00 +0000] "GET / HTTP/2.0" 200 [...]
```

**Apache**:

```
192.0.2.2 - - [13/Dec/1989:13:37:00 +0000] "GET / HTTP/1.1" 302 [...]
```

## How do I install it?
If your system is running Go (usually the `golang` package), then a
`go get github.com/daknob/hstats` and a `go install` should suffice.
This is recommended for systems running Go 1.6 or above since v1.3.3
shipped with Debian is too old and may be slower or not support some
features.

If you don't want to install Go in your servers, which makes total
sense, you can either follow the above instructions on another computer
and then transfer the binaries to the servers, or download the correct
binary for your architecture from the GitHub Releases feature when they
are provided.

## How do I use it?
In order to use `hstats` you currently have to pipe to its `stdin` a
compliant log file. Future versions may or may not support selecting
one or more files, however the easiest way to run this is:

```bash
cat /var/log/nginx/access.log /var/log/nginx/access.log.1 [...] | hstats
```

If you want to examine older files than that, which are compressed and
stored with the `.gz` extension, you can simply use `zcat`, the `cat`
equivalent that first decompresses the files:

```bash
zcat /var/log/nginx/access.log.5 | hstats
```

For more information on how to use the tool, run `hstats -h`.

## What should I expect from it?
A sample output of `hstats` when run against a normal NGiNX log file
can be found below, with several command line arguments passed to
the program:

No command line arguments:

```bash
$ cat access.log | hstats
HTTP/2   Requests:            0 -- 0.00%
HTTP/1.1 Requests:         1419 -- 99.86%
HTTP/1   Requests:            2 -- 0.14%
--
IPv4     Requests:         1421 -- 100.00%
IPv6     Requests:            0 -- 0.00%
--
HTTP 2XX Requests:          924 -- 65.02%
HTTP 3XX Requests:            6 -- 0.42%
HTTP 4XX Requests:          491 -- 34.55%
HTTP 5XX Requests:            0 -- 0.00%
```

With the `-parseline` command line argument:

```bash
$ cat access.log | hstats -parseline
0
1419
2
1421
0
924
6
491
0
```

With the `-parseflat` command line argument:

```bash
$ cat access.log | hstats -parseflat
0 1419 2 1421 0 924 6 491 0
```

With the `-h` command line argument:

```bash
$ hstats -h
Usage of /usr/local/bin/hstats:
  -human=true: Print the output in a human-readable format.
  -parseflat=false: Output statistics for parsing, one line, separated by spaces.
  -parseline=false: Output statistics for parsing, one in each line.
  -showorder=false: Show the order in which the parse output is being printed.
```

With the `-showorder` flag:

```bash
$ hstats -showorder
HTTP2_Requests, HTTP11_Requests, HTTP1_Requests, IPv4_Requests,
IPv6_Requests, HTTP2XX_Requests, HTTP3XX_Requests, HTTP4XX_Requests,
HTTP5XX_Requests
```
