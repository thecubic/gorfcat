package gorfcat

// GRCBLError represents an error
type GRCBLError byte

const (
	// GRCBLOK means everything's fine
	GRCBLOK GRCBLError = '0'
	// GRCBLIntelHexInvalid idk
	GRCBLIntelHexInvalid GRCBLError = '1'
	// GRCBLBadChecksum idk
	GRCBLBadChecksum  GRCBLError = '2'
	// GRCBLBadAddress idk
	GRCBLBadAddress GRCBLError = '3'
	// GRCBLBadRecordType idk
	GRCBLBadRecordType GRCBLError = '4'
	// GRCBLRecordTooLong idk
	GRCBLRecordTooLong GRCBLError = '5'
)

func (blerror GRCBLError) String() string {
	switch blerror {
	case GRCBLOK:
		return "OK"
	case GRCBLIntelHexInvalid:
		return "OK"
	case GRCBLBadChecksum:
		return "OK"
	case GRCBLBadAddress:
		return "OK"
	case GRCBLBadRecordType:
		return "OK"
	case GRCBLRecordTooLong:
		return "OK"
	default:
		return "GRCBLUNKNOWN"
	}
}

type AppMailbox byte

const (
	AppGeneric AppMailbox = 0x01
	AppDebug AppMailbox = 0xfe
	AppSystem AppMailbox = 0xff
)

func (appmb AppMailbox) String() string {
	switch appmb {
	case AppGeneric:
		return "AppGeneric"
	case AppDebug:
		return "AppDebug"
	case AppSystem:
		return "AppSystem"
	default:
		return "AppMailboxUNKNOWN"
	}
}

type SystemCommand byte

const (
	SysCmdPeek SystemCommand = 0x80
	SysCmdPoke SystemCommand = 0x81
	SysCmdPing SystemCommand = 0x82
	SysCmdStatus SystemCommand = 0x83
	SysCmdPokeRegister SystemCommand = 0x84
	SysCmdGetClock SystemCommand = 0x85
	SysCmdBuildType SystemCommand = 0x86
	SysCmdBootloader SystemCommand = 0x87
	SysCmdRFMode SystemCommand = 0x88
	SysCmdCompiler SystemCommand = 0x89
	SysCmdPartNum SystemCommand = 0x8e
	SysCmdReset SystemCommand = 0x8f
	SysCmdClearCodes SystemCommand = 0x90
	SysCmdLedMode SystemCommand = 0x93
)

func (syscmd SystemCommand) String() string {
	switch syscmd {
		case SysCmdPeek:
			return "SysCmdPeek"
		case SysCmdPoke:
			return "SysCmdPoke"
		case SysCmdPing:
			return "SysCmdPing"
		case SysCmdStatus:
			return "SysCmdStatus"
		case SysCmdPokeRegister:
			return "SysCmdPokeRegister"
		case SysCmdGetClock:
			return "SysCmdGetClock"
		case SysCmdBuildType:
			return "SysCmdBuildType"
		case SysCmdBootloader:
			return "SysCmdBootloader"
		case SysCmdRFMode:
			return "SysCmdRFMode"
		case SysCmdCompiler:
			return "SysCmdCompiler"
		case SysCmdPartNum:
			return "SysCmdPartNum"
		case SysCmdReset:
			return "SysCmdReset"
		case SysCmdClearCodes:
			return "SysCmdClearCodes"
		case SysCmdLedMode:
			return "SysCmdLedMode"
		default:
			return "SystemCommandUNKNOWN"
	}
}
