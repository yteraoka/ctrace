package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	trace "cloud.google.com/go/trace/apiv1"
	cloudtracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
)

const (
	EXIT_UNKNOWN = 1
	EXIT_NOT_FOUND = 20
)

var optQuiet bool
var optUrl bool

func init() {
	flag.BoolVar(&optQuiet, "quiet", false, "no output")
	flag.BoolVar(&optQuiet, "q", false, "no output (shorthand)")
	flag.BoolVar(&optUrl, "url", false, "print cloud trace url")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-quiet] [-url] projects/PROJECT_ID/traces/xxxxxxxx\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s [-quiet] [-url] PROJECT_ID xxxxxxxx\n", os.Args[0])
	os.Exit(EXIT_UNKNOWN)
}

func main() {
	flag.Parse()

	var projectId string
	var traceId string

	if len(flag.Args()) == 1 {
		parts := strings.Split(flag.Arg(0), "/")
		if parts[0] == "projects" {
			projectId = parts[1]
		}
		if parts[2] == "traces" {
			traceId = parts[3]
		}
	} else if len(flag.Args()) == 2 {
		projectId = flag.Arg(0)
		traceId = flag.Arg(1)
	} else {
		usage()
	}

	ctx := context.Background()
	c, err := trace.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	req := &cloudtracepb.GetTraceRequest{
		ProjectId: projectId,
		TraceId: traceId,
	}
	resp, err := c.GetTrace(ctx, req)
	st, ok := status.FromError(err)

	if err != nil {
		if !ok {
			if ! optQuiet {
				log.Printf("Unknown Error: %+v\n", err)
			}
			os.Exit(EXIT_UNKNOWN)
		}

		if st.Code() == codes.NotFound {
			if ! optQuiet {
				fmt.Printf("trace %s does not exist\n", traceId)
			}
			os.Exit(EXIT_NOT_FOUND)
		} else {
			if ! optQuiet {
				log.Printf("Unknown Error: %+v\n", err)
			}
			os.Exit(EXIT_UNKNOWN)
		}
	}

	if optUrl {
		fmt.Printf("https://console.cloud.google.com/traces/list?tid=%s&project=%s\n",
			resp.GetTraceId(), projectId)
	} else {
		fmt.Printf("%s\n", resp.GetTraceId())
	}
}
