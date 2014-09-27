package main

import (
	"fmt"
	"strings"

	"github.com/AudriusButkevicius/cli"
)

func init() {
	cliCommands = append(cliCommands, cli.Command{
		Name:     "options",
		HideHelp: true,
		Usage:    "Options command group",
		Subcommands: []cli.Command{
			{
				Name:     "dump",
				Usage:    "Show all Syncthing option settings",
				Requires: &cli.Requires{},
				Action:   optionsDump,
			},
			{
				Name:     "get",
				Usage:    "Get a Syncthing option setting",
				Requires: &cli.Requires{"setting"},
				Action:   optionsGet,
			},
			{
				Name:     "set",
				Usage:    "Set a Syncthing option setting",
				Requires: &cli.Requires{"setting", "value..."},
				Action:   optionsSet,
			},
		},
	})
}

func optionsDump(c *cli.Context) {
	cfg := getConfig(c).Options
	writer := newTableWriter()

	fmt.Fprintln(writer, "Sync protocol listen addresses:\t", strings.Join(cfg.ListenAddress, " "), "\t(address)")
	fmt.Fprintln(writer, "Global discovery enabled:\t", cfg.GlobalAnnEnabled, "\t(globalannenabled)")
	fmt.Fprintln(writer, "Global discovery server:\t", cfg.GlobalAnnServer, "\t(globalannserver)")

	fmt.Fprintln(writer, "Local discovery enabled:\t", cfg.LocalAnnEnabled, "\t(localannenabled)")
	fmt.Fprintln(writer, "Local discovery port:\t", cfg.LocalAnnPort, "\t(localannport)")

	fmt.Fprintln(writer, "Outgoing rate limit in KiB/s:\t", cfg.MaxSendKbps, "\t(maxsend)")
	fmt.Fprintln(writer, "Incoming rate limit in KiB/s:\t", cfg.MaxRecvKbps, "\t(maxrecv)")
	fmt.Fprintln(writer, "Reconnect interval in seconds:\t", cfg.ReconnectIntervalS, "\t(reconnect)")
	fmt.Fprintln(writer, "Start browser:\t", cfg.StartBrowser, "\t(browser)")
	fmt.Fprintln(writer, "Enable UPnP:\t", cfg.UPnPEnabled, "\t(upnp)")
	fmt.Fprintln(writer, "UPnP Lease in minutes:\t", cfg.UPnPLease, "\t(upnplease)")
	fmt.Fprintln(writer, "UPnP Renewal period in minutes:\t", cfg.UPnPRenewal, "\t(upnprenew)")
	fmt.Fprintln(writer, "Restart on Wake Up:\t", cfg.RestartOnWakeup, "\t(wake)")

	reporting := "unrecognized value"
	switch cfg.URAccepted {
	case -1:
		reporting = "false"
	case 0:
		reporting = "undecided/false"
	case 1:
		reporting = "true"
	}
	fmt.Fprintln(writer, "Anonymous usage reporting:\t", reporting, "\t(reporting)")

	writer.Flush()
}

func optionsGet(c *cli.Context) {
	cfg := getConfig(c).Options
	arg := c.Args()[0]
	switch strings.ToLower(arg) {
	case "address":
		fmt.Println(strings.Join(cfg.ListenAddress, "\n"))
	case "globalannenabled":
		fmt.Println(cfg.GlobalAnnEnabled)
	case "globalannserver":
		fmt.Println(cfg.GlobalAnnServer)
	case "localannenabled":
		fmt.Println(cfg.LocalAnnEnabled)
	case "localannport":
		fmt.Println(cfg.LocalAnnPort)
	case "maxsend":
		fmt.Println(cfg.MaxSendKbps)
	case "maxrecv":
		fmt.Println(cfg.MaxRecvKbps)
	case "reconnect":
		fmt.Println(cfg.ReconnectIntervalS)
	case "browser":
		fmt.Println(cfg.StartBrowser)
	case "upnp":
		fmt.Println(cfg.UPnPEnabled)
	case "upnplease":
		fmt.Println(cfg.UPnPLease)
	case "upnprenew":
		fmt.Println(cfg.UPnPRenewal)
	case "reporting":
		switch cfg.URAccepted {
		case -1:
			fmt.Println("false")
		case 0:
			fmt.Println("undecided/false")
		case 1:
			fmt.Println("true")
		default:
			fmt.Println("unknown")
		}
	case "wake":
		fmt.Println(cfg.RestartOnWakeup)
	default:
		die("Invalid setting: " + arg + "\nAvailable settings: address, globalannenabled, globalannserver, localannenabled, localannport, maxsend, maxrecv, reconnect, browser, upnp, upnplease, upnprenew, reporting, wake")
	}
}

func optionsSet(c *cli.Context) {
	config := getConfig(c)
	arg := c.Args()[0]
	val := c.Args()[1]
	switch strings.ToLower(arg) {
	case "address":
		for _, item := range c.Args().Tail() {
			validAddress(item)
		}
		config.Options.ListenAddress = c.Args().Tail()
	case "globalannenabled":
		config.Options.GlobalAnnEnabled = parseBool(val)
	case "globalannserver":
		validAddress(val)
		config.Options.GlobalAnnServer = val
	case "localannenabled":
		config.Options.LocalAnnEnabled = parseBool(val)
	case "localannport":
		config.Options.LocalAnnPort = parsePort(val)
	case "maxsend":
		config.Options.MaxSendKbps = parseUint(val)
	case "maxrecv":
		config.Options.MaxRecvKbps = parseUint(val)
	case "reconnect":
		config.Options.ReconnectIntervalS = parseUint(val)
	case "browser":
		config.Options.StartBrowser = parseBool(val)
	case "upnp":
		config.Options.UPnPEnabled = parseBool(val)
	case "upnplease":
		config.Options.UPnPLease = parseUint(val)
	case "upnprenew":
		config.Options.UPnPRenewal = parseUint(val)
	case "reporting":
		switch strings.ToLower(val) {
		case "u", "undecided", "unset":
			config.Options.URAccepted = 0
		default:
			boolvalue := parseBool(val)
			if boolvalue {
				config.Options.URAccepted = 1
			} else {
				config.Options.URAccepted = -1
			}
		}
	case "wake":
		config.Options.RestartOnWakeup = parseBool(val)
	default:
		die("Invalid setting: " + arg + "\nAvailable settings: address, globalannenabled, globalannserver, localannenabled, localannport, maxsend, maxrecv, reconnect, browser, upnp, upnplease, upnprenew, reporting, wake")
	}
	setConfig(c, config)
}
