package main

import (
	"os"
	"time"

	"github.com/petermeissner/loggo/src/util"
	"github.com/urfave/cli/v2"
)

func main() {

	// cli application
	cli_app := &cli.App{

		// name
		Name:    "loggo-simulate",
		Version: "0.0.1",

		// description
		Usage: "Simulate log file creation by reading files and writing new ones",

		// Default Action
		Action: func(cli_context *cli.Context) error {

			// show main help
			cli.ShowAppHelp(cli_context)

			// show help for read
			println("\n--------------------------------------------------------------------\n")
			cli.ShowCommandHelp(cli_context, "emit")

			return nil
		},

		// Sub commands (actions) available
		Commands: []*cli.Command{

			// Command: version
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version",

				Action: func(cli_context *cli.Context) error {
					cli.ShowVersion(cli_context)

					return nil
				},
			},

			// Command: emit
			{
				Name:  "emit",
				Usage: "development",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "logpath",
						Value: "./logs/test.log",
						Usage: "path to log files to server",
					},

					&cli.StringFlag{
						Name:  "msg_pattern",
						Value: `^\d{4}-\d{2}-\d{2}`,
						Usage: "regular expression to decide if a new log message starts",
					},

					&cli.StringFlag{
						Name:  "ts_pattern",
						Value: `yyyy-MM-dd HH:mm:ss,SSS`,
						Usage: "Log message time stamps are parsed and replaced with current time + log time offset. Time stamp pattern can be specified using yyyy as years, MM as months, dd as days, HH as hours, mm as minutes, ss as seconds, S as sub seconds.",
					},

					&cli.StringFlag{
						Name:  "max_mpm",
						Value: `0`,
						Usage: "Either '0' - default - or an integer indicating how many messages per minute should be emitted. If set to '0' the rate is calculated from the log file used as input.",
					},

					&cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Either '' - default: writing to terminal - or an file path indicating where to write the log messages to.",
					},

					&cli.StringFlag{
						Name:  "max_size",
						Value: "0",
						Usage: "Either '0' - default: no size constraint - until which size messages are emitted.",
					},

					&cli.StringFlag{
						Name:  "max_messages",
						Value: "0",
						Usage: "Either '0' - default - maximum number of messages emitted.",
					},

					&cli.StringFlag{
						Name:  "max_time",
						Value: "0",
						Usage: "Either '0' - default - maximum run time in which messages are emitted.",
					},
				},

				Action: func(cli_context *cli.Context) error {

					// setup logger
					logger := util.Setup_logger("info", "text")

					// get parameters
					logpath := cli_context.String("logpath")
					msg_pattern := cli_context.String("msg_pattern")
					ts_pattern := cli_context.String("ts_pattern")
					max_mpm := cli_context.Float64("max_mpm")
					output := cli_context.String("output")
					max_size := cli_context.Float64("max_size")
					max_messages := cli_context.Int64("max_messages")
					max_time := cli_context.Int64("max_time")
					max_time_duration := time.Duration(max_time * int64(time.Second))

					// handle how messages are emitted - to terminal or to file
					var emit func(string)
					if output == "" {

						// set emit function
						emit = func(message string) {
							println(message)
						}
					} else {

						// open file for writing, ensure it gets closed
						file_out, err := os.Create(output)
						if err != nil {
							panic(err)
						}
						defer file_out.Close()

						// set emit function
						emit = func(message string) {
							file_out.WriteString(message)
						}
					}

					// read and parse log file
					messages, e := read_log_messages(logpath, msg_pattern)
					if e != nil {
						logger.Error(e.Error())
						return e
					}

					// handle max_mpm setting
					// - calculate messages per minute rate found in log file if not set via parameter
					max_mpm = max_mpm_handler(max_mpm, cli_context.IsSet("max_mpm"), ts_pattern, messages)

					// emit messages
					i := int64(0)
					start_time := time.Now()
					size := float64(0)
					mpm := float64(0)

					// log info
					logger.Info("Config: ", "max_messages", max_messages, "max_time", max_time, "max_size", max_size, "max_mpm", max_mpm, "output", output, "logpath", logpath, "msg_pattern", msg_pattern, "ts_pattern", ts_pattern, "max_time_duration", max_time_duration)

				main_loop:
					for {
						// prepare a new set of messages
						messages = update_time_stamp(messages, ts_pattern)

						// write messages to output
						for _, message := range messages {

							// emit message
							emit(message)

							// update loop variables
							i++
							size += float64(len(message))

							if i == 1 || i%100 == 0 {
								logger.Info("Progress: ", "i", i, "size", size, "mpm", mpm, "time", time.Since(start_time).String())
							}

							// stop emitting?
							if max_messages != 0 && i >= max_messages ||
								(max_time != 0 && (time.Since(start_time) >= max_time_duration)) ||
								(max_size != 0 && size >= max_size) {

								break main_loop
							}

							// wait before emitting next message?
							if max_mpm != 0 {
								mpm = float64(i) / (time.Since(start_time).Minutes() + 1)
								if mpm >= max_mpm {
									overrate := mpm - max_mpm
									waiting_time := time.Duration(overrate * float64(time.Minute) / 60)
									time.Sleep(waiting_time)
								}
							}
						}

					}

					// done
					return nil
				},
			},
		},
	}

	if err := cli_app.Run(os.Args); err != nil {
		// setup logger
		logger := util.Setup_logger("info", "text")

		logger.Error(err.Error())
	}
}
