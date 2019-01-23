package gorfcat

import (
	"github.com/google/gousb"
)

var (
	TIVendor = gousb.ID(0x0451)
	TIProduct = gousb.ID(0x4715)
	OpenMokoVendor = gousb.ID(0x1d50)
	RFCatProduct1 = gousb.ID(0x6047)
	RFCatProduct2 = gousb.ID(0x6048)
	YardStickOneProduct = gousb.ID(0x605b)
	YardStickOneBootloaderProduct = gousb.ID(0x605c)
	PandwaRFProduct = gousb.ID(0x60ff)
)


func RFCatFilter(desc *gousb.DeviceDesc) bool {
	switch desc.Vendor {
	case TIVendor:
		// TI USB Classic
		return desc.Product == TIProduct
	case OpenMokoVendor:
		// OpenMoko
		switch desc.Product {
			case RFCatProduct1:
				return true
			case RFCatProduct2:
				return true
			case YardStickOneProduct:
				// Yard Stick One
				return true
			case YardStickOneBootloaderProduct:
				// Yard Stick One, but in bootloader mode
				return false
			case PandwaRFProduct:
				return true
			default:
				return false
		}
	default:
		return false
	}
}

// func RFCatBootloaderFilter
