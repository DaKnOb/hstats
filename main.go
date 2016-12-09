package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	/* Open stdin to read logs from */
	logInput, err := os.Open("/dev/stdin")

	/* If there's a problem, log it and exit */
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	/* Close the "file" after we're done */
	defer logInput.Close()

	/* HTTP Protocol Version Used */
	H2 := 0
	H11 := 0
	H1 := 0

	/* IP Protocol Version Used */
	IPv4 := 0
	IPv6 := 0

	/* HTTP Status Code Returned */
	R2XX := 0
	R3XX := 0
	R4XX := 0
	R5XX := 0

	/* Start a scanner to parse the file */
	scanner := bufio.NewScanner(logInput)

	/* Parse every line */
	for scanner.Scan() {
		f := scanner.Text()

		/* Look for HTTP Protocol Version in the log lines */
		if strings.Contains(f, "HTTP/2") {
			H2++
		} else if strings.Contains(f, "HTTP/1.1") {
			H11++
		} else if strings.Contains(f, "HTTP/1.0") {
			H1++
		}

		/*
			Check the first part of each log line ( IP Address )
			and attempt to determine if it's an IPv4 or IPv6 address.
			The IPv4 check must go first to handle IPv4 over tcp6
			sockets (::ffff:192.0.2.2). We currently don't handle
			cases where the tcp6 socket fully translates the IPv4
			address to hex (::ffff:c000:202).
		*/
		if strings.Contains(strings.Split(f, " ")[0], ".") {
			IPv4++
		} else if strings.Contains(strings.Split(f, " ")[0], ":") {
			IPv6++
		}

		/*
			Check the specific position in the NGiNX default log
			file where the HTTP Status Code normally appears.
		*/
		if strings.HasPrefix(strings.Replace(strings.Split(f, "\"")[2], " ", "", -1), "2") {
			R2XX++
		} else if strings.HasPrefix(strings.Replace(strings.Split(f, "\"")[2], " ", "", -1), "3") {
			R3XX++
		} else if strings.HasPrefix(strings.Replace(strings.Split(f, "\"")[2], " ", "", -1), "4") {
			R4XX++
		} else if strings.HasPrefix(strings.Replace(strings.Split(f, "\"")[2], " ", "", -1), "5") {
			R5XX++
		}

	}

	/* Handle scanner initialization errors */
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}

	/* Calculate Percentages */
	TotalHTTPV := float64(H2 + H11 + H1)
	TotalIPV := float64(IPv4 + IPv6)
	TotalRS := float64(R2XX + R3XX + R4XX + R5XX)

	var H2P, H11P, H1P, v4P, v6P, R2P, R3P, R4P, R5P float64

	if TotalHTTPV != 0 {
		H2P = 100 * float64(H2) / TotalHTTPV
		H11P = 100 * float64(H11) / TotalHTTPV
		H1P = 100 * float64(H1) / TotalHTTPV
	}

	if TotalIPV != 0 {
		v4P = 100 * float64(IPv4) / TotalIPV
		v6P = 100 * float64(IPv6) / TotalIPV
	}

	if TotalRS != 0 {
		R2P = 100 * float64(R2XX) / TotalRS
		R3P = 100 * float64(R3XX) / TotalRS
		R4P = 100 * float64(R4XX) / TotalRS
		R5P = 100 * float64(R5XX) / TotalRS
	}

	/* Check which format the output is needed to be printed */
	orderPrint := flag.Bool("showorder", false, "Show the order in which the parse output is being printed.")
	parseLines := flag.Bool("parseline", false, "Output statistics for parsing, one in each line.")
	parseFlat := flag.Bool("parseflat", false, "Output statistics for parsing, one line, separated by spaces.")
	humanReadable := flag.Bool("human", true, "Print the output in a human-readable format.")

	/* Parse the command line flags */
	flag.Parse()

	/* Prints the order in which parseLines and parseFlat numbers are */
	if *orderPrint == true {
		fmt.Println("HTTP2_Requests, HTTP11_Requests, HTTP1_Requests, IPv4_Requests,")
		fmt.Println("IPv6_Requests, HTTP2XX_Requests, HTTP3XX_Requests, HTTP4XX_Requests,")
		fmt.Println("HTTP5XX_Requests")
		os.Exit(0)
	}

	/* Prints the statistics, one per line */
	if *parseLines == true {
		fmt.Printf("%d\n%d\n%d\n%d\n%d\n%d\n%d\n%d\n%d\n", H2, H11, H1, IPv4, IPv6, R2XX, R3XX, R4XX, R5XX)
		os.Exit(0)
	}

	/* Prints the statistics, space delimited */
	if *parseFlat == true {
		fmt.Println(H2, H11, H1, IPv4, IPv6, R2XX, R3XX, R4XX, R5XX)
		os.Exit(0)
	}

	/* Prints the statistics in a human readable format */
	if *humanReadable == true {
		if TotalHTTPV != 0 {
			fmt.Printf("HTTP/2   Requests: %12d -- %3.2f%%\n", H2, H2P)
			fmt.Printf("HTTP/1.1 Requests: %12d -- %3.2f%%\n", H11, H11P)
			fmt.Printf("HTTP/1   Requests: %12d -- %3.2f%%\n", H1, H1P)
			fmt.Printf("--\n")
		}
		if TotalIPV != 0 {
			fmt.Printf("IPv4     Requests: %12d -- %3.2f%%\n", IPv4, v4P)
			fmt.Printf("IPv6     Requests: %12d -- %3.2f%%\n", IPv6, v6P)
			fmt.Printf("--\n")
		}
		if TotalRS != 0 {
			fmt.Printf("HTTP 2XX Requests: %12d -- %3.2f%%\n", R2XX, R2P)
			fmt.Printf("HTTP 3XX Requests: %12d -- %3.2f%%\n", R3XX, R3P)
			fmt.Printf("HTTP 4XX Requests: %12d -- %3.2f%%\n", R4XX, R4P)
			fmt.Printf("HTTP 5XX Requests: %12d -- %3.2f%%\n", R5XX, R5P)
		}
		os.Exit(0)
	}

	/* This is added here just in case */
	fmt.Println("Please consider the -help manual for this tool.")
	fmt.Println("You need to specify the output format for this to work.")
	os.Exit(3)

}
