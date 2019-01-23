package gorfcat

import (
	"fmt"
	"io"
	"bytes"
	"os"
	"syscall"
	"github.com/marcinbor85/gohex"
	// "go.bug.st/serial.v1"
)

// OpenDevice opens a bootloader device with appropriate settings
// ACM is a UART abstraction, so it is just a regular file
func OpenDevice(device string) (*os.File, error) {
	// don't be a controlling terminal (needed?)
	return os.OpenFile(device, os.O_RDWR|syscall.O_NOCTTY, 0)
}

func HexLoadFile(fd io.Reader) (*gohex.Memory, error) {
	mem := gohex.NewMemory()
	err := mem.ParseIntelHex(fd)
	return mem, err
}

func HexLoad(blob []byte) (*gohex.Memory, error) {
	mem := gohex.NewMemory()
	err := mem.ParseIntelHex(bytes.NewReader(blob))
	return mem, err
}

// RunUserCode is a one-way function that exits
// the bootloader (it's also an EOF record)
func RunUserCode(device *os.File) error {
	_, err := device.Write([]byte(":00000001FF\n"))
	return err
}

// ain't work because Read ain't work
func Verify(device *os.File, ihex *gohex.Memory) error {
	var err error
	datasegments := ihex.GetDataSegments()
	for _, segment := range datasegments {
		fmt.Printf("Segment: @%x %d %x\n", segment.Address, len(segment.Data), segment.Data)
		err = Read(device, uint16(segment.Address), uint16(len(segment.Data)))
		if err != nil {
			return err
		}
	}
	return err
}

// NOTE: ain't work yet
func Read(device *os.File, address uint16, length uint16) error {
	var (
		err error
		buf []byte
		bufi int
	)
	// weird flex but ok
	chksum := (0xD9 +
			   (0x100 - (address & 0xFF)) +
			   (0x100 - ((address >> 8) & 0xFF)) +
			   (0x100 - (length & 0xFF)) +
			   (0x100 - ((length >> 8) & 0xFF))) & 0xFF
	// :02[address]25[length][checksum]
	strcmd := fmt.Sprintf(":02%04X25%04X%02X\n", address, length, chksum)
	fmt.Printf("strcmd: %v", strcmd)
	device.Write([]byte(strcmd))
	for {
		bufi, err = device.Read(buf)
		fmt.Printf("read %d bytes\n", bufi)
		fmt.Printf("read: %v\n", buf[:bufi])
		fmt.Printf("err: %v\n\n", err)
		// don't think this can happen
		if err == io.EOF {
			fmt.Println("eof")
			break
		} else if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
		// } else if line == "\n" {
		// 	fmt.Println("finished")
		// 	break
		// } else {
		// 	fmt.Printf("line: ", line)
		// }
	}
	return err
}

// ResetPageClearProtection -> ":00000022DE\n"
// Reset record will reset the page erase map which usually ensures each page is only
// erased once, allowing for random writes but preventing overwriting of data already written
// this session.
// <- rc
func ResetPageClearProtection(device *os.File) error {
	rc := make([]byte, 1)
	_, err := device.Write([]byte(":00000022DE\n"))
	if err != nil {
		return err
	}
	_, err = device.Read(rc)
	if err != nil {
		return err
	} else if GRCBLError(rc[0]) == GRCBLOK {
		return nil
	} else {
		return fmt.Errorf("bad response: %v", GRCBLError(rc[0]))
	}
}

// EraseAllUserCode erases all user code flash pages
// -> ":00000023DD\n"
// <- rc

// EraseUserPage(page) <- record 24 page number and checksum

