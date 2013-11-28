// Package telnet,
// Implements RFC 854
// Features:
//     - Implements client related functions
//     - Server related
//     - Generic io.Reader and io.Writer to connect to any other programs

//  RFC 854  Telnet Protocol Specification
//  RFC 855  Telnet Option Specifications
//  RFC 856  Telnet Binary Transmission
//  RFC 857  Telnet Echo Option
//  RFC 858  Telnet Suppress Go Ahead option
//  RFC 859  Telnet Status Option
//  RFC 860  Telnet Timing Mark option
//  RFC 861  Telnet Extended Options option
//  RFC 1091 Telnet Terminal Type option
//  RFC 1096 Telnet X Display location option
//  RFC 1184 Telnet Linemode option
//  RFC 1372 Telnet Remote Flow option
//  RFC 1408 Telnet Environment option
//  RFC 1571 Telnet Environment option interoperability issues
//  RFC 1572 Telnet Environment option
//  RFC 2066 Telnet Charset option
//  RFC 2941 Telnet Authentication option
//  RFC 2840 Telnet Kermit option
//  RFC 2217 Telnet Com Port option
//  RFC 1073 Telnet Window Size option
//  RFC 1079 Telnet Terminal Speed option
//  RFC 727  Telnet logout option

package telnet

import (
	"bufio"
	//"bytes"
	"fmt"
	//"io"
	//"log"
	"net"
	//"strings"
	"time"
)

// Each sequence of read is one or two bytes, depending on if the byte read
// has value {0-254} or 255, if later(255) case, then its a control command,
// and next byte indicates the actual command

type Telnet struct {
	conn     net.Conn
	bufrd    *bufio.Reader
	bufwr    *bufio.Writer // Dont know if buffered writer required
	unixCRLF bool

	debug bool
	cmd   byte // Previous CMD
	opt   byte // Option for previous CMD
}

func (t *Telnet) EnableDebug() {
	t.debug = true
}

type TelnetServer struct {
	Telnet
	listn net.Listener
}

const (
	CR  = byte('\r')
	LF  = byte('\n')
	NUL = byte(0)
)

type telnetCMD byte

const (
	cmd_EOF   = 236 + iota // End of file
	cmd_SUSP               // 237: Suspend process
	cmd_ABORT              // Abort process
	cmd_EOR                // end of record
	cmd_SE                 // 240: end of sub negotiation
	cmd_NOP                // 241: nop
	cmd_DM                 // data mark
	cmd_BREAK              // break
	cmd_IP                 // interrupt process
	cmd_AO                 // abort output
	cmd_AYT                // 246: Are You There
	cmd_EC                 // delete current character
	cmd_EL                 // delete current line
	cmd_GA                 // Line reverse
	cmd_SB                 // 250: Subnegotiation
	cmd_WILL               // Indicating Option *WILL* be used
	cmd_WONT               // Indicating option *WONT* be used
	cmd_DO                 // Commanding, to use Option
	cmd_DONT               // Response, Option not supported
	cmd_IAC                // 255: Interpret As Command
)

var cmdStrings = [256]string{
	// 0 - 236 are not used
	cmd_EOF:   "EOF",
	cmd_SUSP:  "SUSP",
	cmd_ABORT: "ABORT",
	cmd_EOR:   "EOR",
	cmd_SE:    "SE",
	cmd_NOP:   "NOP",
	cmd_DM:    "DM",
	cmd_BREAK: "BREAK",
	cmd_IP:    "IP",
	cmd_AO:    "AO",
	cmd_AYT:   "AYT",
	cmd_EC:    "EC",
	cmd_EL:    "EL",
	cmd_GA:    "GA",
	cmd_SB:    "SB",
	cmd_WILL:  "WILL",
	cmd_WONT:  "WONT",
	cmd_DO:    "DO",
	cmd_DONT:  "DONT",
	cmd_IAC:   "IAC",
}

func (c telnetCMD) String() string {
	return cmdStrings[c]
}

type telnetOPT byte

const (
	opt_BINARY         = iota // 8-bit data path
	opt_ECHO                  // echo
	opt_RCP                   // prepare to reconnect
	opt_SGA                   // suppress go ahead
	opt_NAMS                  // approximate message size
	opt_STATUS                // give status
	opt_TM                    // timing mark
	opt_RCTE                  // remote controlled transmission and echo
	opt_NAOL                  // negotiate about output line width
	opt_NAOP                  // negotiate about output page size
	opt_NAOCRD                // negotiate about CR disposition
	opt_NAOHTS                // negotiate about horizontal tabstops
	opt_NAOHTD                // negotiate about horizontal tab disposition
	opt_NAOFFD                // negotiate about formfeed disposition
	opt_NAOVTS                // negotiate about vertical tab stops
	opt_NAOVTD                // negotiate about vertical tab disposition
	opt_NAOLFD                // negotiate about output LF disposition
	opt_XASCII                // extended ascic character set
	opt_LOGOUT                // force logout
	opt_BM                    // byte macro
	opt_DET                   // data entry terminal
	opt_SUPDUP                // supdup protocol
	opt_SUPDUPOUTPUT          // supdup output
	opt_SNDLOC                // send location
	opt_TTYPE                 // terminal type
	opt_EOR                   // end or record
	opt_TUID                  // TACACS user identification
	opt_OUTMRK                // output marking
	opt_TTYLOC                // terminal location number
	opt_3270REGIME            // 3270 regime
	opt_X3PAD                 // X.3 PAD
	opt_NAWS                  // window size
	opt_TSPEED                // terminal speed
	opt_LFLOW                 // remote flow control
	opt_LINEMODE              // Linemode option
	opt_XDISPLOC              // X Display Location
	opt_OLD_ENVIRON           // Old - Environment variables
	opt_AUTHENTICATION        // Authenticate
	opt_ENCRYPT               // Encryption option
	opt_NEW_ENVIRON           // New - Environment variables
	opt_EXOPL          = 255  // extended-options-list
)

var optStrings = [256]string{
	opt_BINARY:         "BINARY",
	opt_ECHO:           "ECHO",
	opt_RCP:            "RCP",
	opt_SGA:            "SGA",
	opt_NAMS:           "NAMS",
	opt_STATUS:         "STATUS",
	opt_TM:             "TM",
	opt_RCTE:           "RCTE",
	opt_NAOL:           "NAOL",
	opt_NAOP:           "NAOP",
	opt_NAOCRD:         "NAOCRD",
	opt_NAOHTS:         "NAOHTS",
	opt_NAOHTD:         "NAOHTD",
	opt_NAOFFD:         "NAOFFD",
	opt_NAOVTS:         "NAOVTS",
	opt_NAOVTD:         "NAOVTD",
	opt_NAOLFD:         "NAOLFD",
	opt_XASCII:         "XASCII",
	opt_LOGOUT:         "LOGOUT",
	opt_BM:             "BM",
	opt_DET:            "DET",
	opt_SUPDUP:         "SUPDUP",
	opt_SUPDUPOUTPUT:   "SUPDUPOUTPUT",
	opt_SNDLOC:         "SNDLOC",
	opt_TTYPE:          "TTYPE",
	opt_EOR:            "EOR",
	opt_TUID:           "TUID",
	opt_OUTMRK:         "OUTMRK",
	opt_TTYLOC:         "TTYLOC",
	opt_3270REGIME:     "REGIME",
	opt_X3PAD:          "X3PAD",
	opt_NAWS:           "NAWS",
	opt_TSPEED:         "TSPEED",
	opt_LFLOW:          "LFLOW",
	opt_LINEMODE:       "LINEMODE",
	opt_XDISPLOC:       "XDISPLOC",
	opt_OLD_ENVIRON:    "OLD_VER",
	opt_AUTHENTICATION: "AUTHENTICATION",
	opt_ENCRYPT:        "ENCRYPT",
	opt_NEW_ENVIRON:    "NEW_ENVIRON",
	// 255
	opt_EXOPL: "EXOPL",
}

func (c telnetOPT) String() string {
	return optStrings[c]
}

//
//
// ######    ##    #####      ##
//   ##    ##  ##  ##  ##   ##  ##
//   ##    ##  ##  ##   ##  ##  ##
//   ##    ##  ##  ##  ##   ##  ##
//   ##      ##    #####      ##
//  Some options have sub-options,
//     opt_AUTHENTICATION, opt_ENCRYPT,
//  And will be implemented in course of time.
//

func NewTelnet() *Telnet {
	return &Telnet{
		unixCRLF: true,
	}
}

func NewTelnetServer() *TelnetServer {
	return &TelnetServer{}
}

func connect(c chan error, t *TelnetServer) {
	var e error
	if t.conn, e = t.listn.Accept(); e != nil {
		c <- e
		return
	}

	t.bufwr = bufio.NewWriterSize(t.conn, 512)
	t.bufrd = bufio.NewReaderSize(t.conn, 512)

	c <- nil
}

// Options are passed like telnet=tcp!localhost:2030
// Change is to accept everything that golang/pkg/net can do with
// 'proto' and 'addr'
// eg:
//        Dial("tcp", "12.34.56.78:80")      OR  Dial("tcp", "google.com:http")
//        Dial("tcp", "[2001:db8::1]:http")  OR  Dial("tcp", "[fe80::1%lo0]:80")
//        Dial("ip4:1", "127.0.0.1")         OR  Dial("ip6:ospf", "::1")
func (t *TelnetServer) ListenTimeout(proto, addr string, dur time.Duration) (e error) {
	if t.listn, e = net.Listen(proto, addr); e != nil {
		return
	}

	con_ch := make(chan error)
	go connect(con_ch, t)

	select {
	case <-time.After(dur * time.Second):
	case e = <-con_ch:
	}

	if e != nil {
		fmt.Printf("%v", e)
	}

	if t.debug {
		if t.conn != nil {
			fmt.Printf("Connected: ")
		}
		fmt.Printf("%v\n", t.conn.LocalAddr())
	}
	return
}

// TODO: FIX THIS, STILL HAS PROBLEMS
// This is similar to ListenTimeout, except, it prints number of seconds waited,
// And exits if Errors are supposed to be treated seriously
func (t *TelnetServer) ListenTimeoutProgress(proto, addr string, dur time.Duration) (e error) {
	timeout := int(dur)
	for ; timeout > 0; timeout-- {
		if t.debug {
			fmt.Printf("Waiting %d seconds for connection \n", timeout)
		}
		if e = t.ListenTimeout(proto, addr, 1); e != nil {
			return
		}
		// In case we are connected we are out
		if t.conn != nil {
			return
		}
	}

	return
}

// io.Writer interface
func (t *Telnet) Write(buf []byte) (n int, err error) {
	for len(buf) > 0 {
		if n, err = t.conn.Write(buf); err != nil {
			return
		}
		buf = buf[n:]
	}
	return
}

// io.Reader interface
func (t *Telnet) Read(buf []byte) (n int, err error) {
	buflen := len(buf)
	for n < buflen {
		if b, err := t.conn.Read(buf); err != nil {
			return b, err
		} else {
			n += b
			buf = buf[b:]
		}
	}
	return n, nil
}

func printcmd(a byte) {
	fmt.Printf("Command received %v\n", telnetCMD(a))
}
func printopt(o byte) {
	fmt.Println("Option received %v\n", telnetOPT(o))
}

func (t *Telnet) __execCMD(c byte) (err error) {

	switch c {
	case cmd_ABORT:
	case cmd_SUSP:
	case cmd_EOR:
	case cmd_NOP:
	case cmd_DM:
	case cmd_BREAK:
		fmt.Println("Sending Backspace")
		t.conn.Write([]byte{cmd_IAC, cmd_EC, cmd_DO})

	case cmd_IP:
		fmt.Println("Sending Backspace")
		t.conn.Write([]byte{cmd_IAC, cmd_EC, cmd_DO})
	case cmd_AO:
	case cmd_AYT:
	case cmd_EC:
	case cmd_EL:
	case cmd_GA:
	case cmd_SE:
	case cmd_SB:
	}

	return
}

func (t *Telnet) _do() {
	if t.debug {
		fmt.Printf("Sending DO %v\n", telnetCMD(t.cmd))
	}
	t.conn.Write([]byte{cmd_IAC, t.cmd, cmd_IAC, cmd_DO})
}

func (t *Telnet) _dont() {
	if t.debug {
		fmt.Printf("Sending DONT %v\n", telnetCMD(t.cmd))
	}
	t.conn.Write([]byte{cmd_IAC, t.cmd, cmd_IAC, cmd_DONT})

}

func (t *Telnet) _will() {
	if t.debug {
		fmt.Printf("Sending WILL %v\n", telnetCMD(t.cmd))
	}
	t.conn.Write([]byte{cmd_IAC, t.cmd, cmd_IAC, cmd_WILL})

}

func (t *Telnet) _wont() {
	if t.debug {
		fmt.Printf("Sending WONT %v\n", telnetCMD(t.cmd))
	}
	t.conn.Write([]byte{cmd_IAC, t.cmd, cmd_IAC, cmd_WONT})
}

func (t *Telnet) __readByte() (c byte, again bool, err error) {
	// We have to interpret the 'telnet' commands and options
	// Send the left overs to whoever asking

	if c, err = t.bufrd.ReadByte(); err != nil {
		return 0, false, err
	}

	if c == cmd_IAC {
		printcmd(c)
		if c, err = t.bufrd.ReadByte(); err != nil {
			return 0, false, err
		}
		printcmd(c)
		switch c {
		case cmd_WILL:
			t._do()
		case cmd_WONT:
			t._dont()
		case cmd_DO:
			t._will()
		case cmd_DONT:
			t._wont()
		default:
			again = true
			t.cmd = c
			err = t.__execCMD(c)
		}
		c = 0
	}

	return
}

// bufio.Reader
func (t *Telnet) ReadByte() (c byte, err error) {
	again := true
	for again == true {
		if c, again, err = t.__readByte(); err != nil {
			c = 0
		}
	}
	return
}

func (t *Telnet) ReadBytes(delim byte) (line []byte, err error) {
	// TODO: need to take care of interpreting the commands
	if delim == 0 {
		delim = LF
	}

	for {
		c, err := t.ReadByte()
		if err != nil {
			break
		}
		line = append(line, c)

		if c == delim {
			break
		}
	}
	return
}

func (t *Telnet) ReadLine() (line []byte, isPrefix bool, err error) {
	return t.bufrd.ReadLine()
}

// func (t *Telnet) ReadRune() (r rune, size int, err error)           {}
// func (t *Telnet) ReadSlice(delim byte) (line []byte, err error)     {}
// func (t *Telnet) ReadString(delim byte) (line string, err error)    {}

func (t *Telnet) Close() {
	if t.debug {
		fmt.Printf("Closing conn:%v\n", t.conn)
	}
	t.conn.Close()
}

func (ts *TelnetServer) Close() {
	//if ts.debug {
	//	fmt.Printf("Closing listn:%v conn:%v\n", ts.listn, ts.conn)
	//}
	ts.Telnet.Close()
	ts.listn.Close()
}
