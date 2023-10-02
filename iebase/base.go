package iebase

func IsOkTidAddr(tid_addr uint16) bool {
	if tid_addr > 0xffff && tid_addr&0xC000 == 0xC000 {
		return false
	}
	return true
}

func IsTidType24(tid_type uint16) bool {
	return tid_type >= 0x24
}
