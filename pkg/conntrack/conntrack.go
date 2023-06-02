package conntrack

import (
	"log"
	"strconv"

	ct "github.com/florianl/go-conntrack"
)

type ConntrackInfo struct {
	SrcAddress   string
	SrcPort      string
	DstAddress   string
	DstPort      string
	TcpStateCode uint8
	TcpState     string
}

type ConntrackInfos struct {
	ConntrackInfos []ConntrackInfo
}

func NewConntrackInfos() *ConntrackInfos {
	return &ConntrackInfos{}
}

func (c *ConntrackInfos) GetConntrack() {
	nfct, err := ct.Open(&ct.Config{})
	if err != nil {
		log.Println("Failed GetConntrack Could not create nfct:", err)
		return
	}
	defer nfct.Close()
	sessions, err := nfct.Dump(ct.Conntrack, ct.IPv4)
	if err != nil {
		log.Println("Failed GetConntrack Could not dump sessions:", err)
		return
	}
	for _, session := range sessions {
		if session.Origin == nil || session.ProtoInfo == nil {
			continue
		}
		log.Println(session.Origin.Src, *session.Origin.Proto.SrcPort, session.Origin.Dst, *session.Origin.Proto.DstPort, *session.ProtoInfo.TCP.State)
		c.ConntrackInfos = append(c.ConntrackInfos, ConntrackInfo{
			SrcAddress:   session.Origin.Src.String(),
			SrcPort:      strconv.Itoa(int(*session.Origin.Proto.SrcPort)),
			DstAddress:   session.Origin.Dst.String(),
			DstPort:      strconv.Itoa(int(*session.Origin.Proto.DstPort)),
			TcpStateCode: *session.ProtoInfo.TCP.State,
			TcpState:     GetStateName(*session.ProtoInfo.TCP.State),
		})
	}
}

func GetStateName(n uint8) string {
	switch n {
	case 0:
		return "LISTEN"
	case 1:
		return "SYN-SENT"
	case 2:
		return "SYN-RECEIVED"
	case 3:
		return "ESTABLISHED"
	case 4:
		return "FIN-WAIT-1"
	case 5:
		return "FIN-WAIT-2"
	case 6:
		return "CLOSE-WAIT"
	case 7:
		return "CLOSING"
	case 8:
		return "LAST-ACK"
	case 9:
		return "TIME-WAIT"
	case 10:
		return "CLOSED"
	}
	return ""
}
