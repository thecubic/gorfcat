package main

import (
	"fmt"
	"os"
	"flag"
	"go.bug.st/serial.v1/enumerator"
	log "github.com/sirupsen/logrus"
	"github.com/thecubic/gorfcat"
	"github.com/google/gousb"
	// "encoding/binary"
)

var (
	debug = flag.Bool("debug", false, "enable debugging messages")
)

func usageQuit() {
	fmt.Println("usage: grc-bootloader <command> [args]")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	command := flag.Arg(0)
	if command == "" {
		usageQuit()
	}

	// device := flag.Arg(0)
	// if device == "" {
	// 	usageQuit()
	// }

	if command == "parse" {
		firmware := flag.Arg(1)
		fwfile, err := os.Open(firmware)
		if err != nil {
			panic(err)
		}
		defer fwfile.Close()
		mem, err := gorfcat.HexLoadFile(fwfile)
		fmt.Printf("mem: %v\n", mem)
		startaddr, ok := mem.GetStartAddress()
		if ok {
			fmt.Printf("Start Address: %x\n", startaddr)
		} else {
			fmt.Println("Start Address: not present")
		}

		datasegments := mem.GetDataSegments()
		for _, segment := range datasegments {
			fmt.Printf("Segment: @%x %d %x\n", segment.Address, len(segment.Data), segment.Data)
		}
	} else if command == "verify" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := gorfcat.OpenDevice(devicefile)
		if err != nil {
			panic(err)
		}
		defer device.Close()

		firmwarefile := flag.Arg(2)
		fwdump, err := os.Open(firmwarefile)
		if err != nil {
			panic(err)
		}
		defer fwdump.Close()
		mem, err := gorfcat.HexLoadFile(fwdump)
		if err != nil {
			panic(err)
		}

		verify := gorfcat.Verify(device, mem)
		fmt.Printf("verify: %v\n", verify)
	} else if command == "list2" {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			fmt.Println("no ports found")
			return
		}
		for _, port := range ports {
			fmt.Printf("Port: %s\n", port.Name)
			if port.IsUSB {
				fmt.Printf("  ID: V%s,P%s\n", port.VID, port.PID)
				fmt.Printf("  Serial: %s\n", port.SerialNumber)
			}
		}
	} else if command == "run" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := gorfcat.OpenDevice(devicefile)
		if err != nil {
			panic(err)
		}
		defer device.Close()
		gorfcat.RunUserCode(device)
	} else if command == "rbcp" {
		devicefile := flag.Arg(1)
		if devicefile == "" {
			usageQuit()
		}
		device, err := gorfcat.OpenDevice(devicefile)
		if err != nil {
			panic(err)
		}
		defer device.Close()
		err = gorfcat.ResetPageClearProtection(device)
		if err != nil {
			panic(err)
		}
	} else if command == "list" {
		usbctx := gousb.NewContext()
		defer usbctx.Close()

		devices, err := usbctx.OpenDevices(gorfcat.RFCatFilter)
		for _, device := range devices {
			defer device.Close()
		}
		if err != nil {
			log.Fatalf("OpenDevices(): %v", err)
		}
		if len(devices) == 0 {
			log.Info("No Devices Found")
		}

		for _, device := range devices {
			// configuration 1
			config, err := device.Config(1)
			if err != nil {
		    	log.Fatalf("%s.Config(1): %v", device, err)				
			}
			defer config.Close()

			// interface 0 altsetting 0
			intf, err := config.Interface(0, 0)
			if err != nil {
    			log.Fatalf("%s.Interface(0, 0): %v", config, err)
			}
			defer intf.Close()

			// In this interface open endpoint #5 for reading.
			epIn, err := intf.InEndpoint(5)
			if err != nil {
				log.Fatalf("%s.InEndpoint(5): %v", intf, err)
			}

			// And in the same interface open endpoint #5 for writing.
			epOut, err := intf.OutEndpoint(5)
			if err != nil {
				log.Fatalf("%s.OutEndpoint(5): %v", intf, err)
			}

			readp := make([]byte, epIn.Desc.MaxPacketSize)

			epOut.Write([]byte{
				byte(gorfcat.AppSystem),
				byte(gorfcat.SysCmdBuildType),
				0x00, 0x00,
			})

			for {
				rn, err := epIn.Read(readp)
				fmt.Printf("read %v\nerr: %v\nbytes: %v\n", rn, err, readp)
				// first byte for Jesus
				mailbox := gorfcat.AppMailbox(readp[1])
				cmd := gorfcat.SystemCommand(readp[2])
				fmt.Printf("mailbox: %v, cmd: %v\n", mailbox, cmd)
				fmt.Printf("payload: %v\n", readp[3:])
				fmt.Printf("string: %v\n", string(readp[5:rn]))
				if cmd == gorfcat.SysCmdBuildType {
					break
				}
			}
		}
	} else if command == "bootloader-all" {
		usbctx := gousb.NewContext()
		defer usbctx.Close()

		devices, err := usbctx.OpenDevices(gorfcat.RFCatFilter)
		for _, device := range devices {
			defer device.Close()
		}
		if err != nil {
			log.Fatalf("OpenDevices(): %v", err)
		}
		if len(devices) == 0 {
			log.Info("No Devices Found")
		}

		for _, device := range devices {
			// configuration 1
			config, err := device.Config(1)
			if err != nil {
				log.Fatalf("%s.Config(1): %v", device, err)				
			}
			defer config.Close()

			// interface 0 altsetting 0
			intf, err := config.Interface(0, 0)
			if err != nil {
				log.Fatalf("%s.Interface(0, 0): %v", config, err)
			}
			defer intf.Close()

			// In this interface open endpoint #5 for reading.
			epIn, err := intf.InEndpoint(5)
			if err != nil {
				log.Fatalf("%s.InEndpoint(5): %v", intf, err)
			}

			// And in the same interface open endpoint #5 for writing.
			epOut, err := intf.OutEndpoint(5)
			if err != nil {
				log.Fatalf("%s.OutEndpoint(5): %v", intf, err)
			}

			readp := make([]byte, epIn.Desc.MaxPacketSize)

			epOut.Write([]byte{
				byte(gorfcat.AppSystem),
				byte(gorfcat.SysCmdBootloader),
				0x00, 0x00,
			})

			for {
				rn, err := epIn.Read(readp)
				fmt.Printf("read %v\nerr: %v\nbytes: %v\n", rn, err, readp)
				// first byte for Jesus
				mailbox := gorfcat.AppMailbox(readp[1])
				cmd := gorfcat.SystemCommand(readp[2])
				fmt.Printf("mailbox: %v, cmd: %v\n", mailbox, cmd)
				fmt.Printf("payload: %v\n", readp[3:])
				fmt.Printf("string: %v\n", string(readp[5:rn]))
				if cmd == gorfcat.SysCmdBootloader {
					break
				}
			}
		}
	}
}
