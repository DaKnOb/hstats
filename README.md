# hstats
A simple NGiNX HTTP Log Parser and Analyzer

## What's that?
`hstats` is a tool written in Go to parse (mostly custom) HTTP Request
Logs from one or more NGiNX web servers. After the parsing is complete,
the tool will output statistics about the requests served by the server
in either a human readable format or a parseable one, for further use in
scripts.

The initial version of the tool was created in 5-10 minutes in order to
help generate statistics really quick, without using any complex software.
It accepts logs that match some certain criteria from `stdin` and then
serially parses them to generate a report in the end.

Future versions of this tool should be able to find the necessary information
in any given log line and possibly be NGiNX variable aware. Until then,
all Pull Requests and Issues are welcome.
