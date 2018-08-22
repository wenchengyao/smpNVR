package cmd

//	"fmt"
//	"net"

type Icmd struct {
	oc     int     //open or close
	op     byte    //operate:1 直播 2点播 3截图 4四宫格 5九宫格 6 六宫格
	cqueue string  //线路
	name   string  //播放名称
	zb     LiveCmd //直播
	db     DbCmd   //点播
	jt     JtCmd   //截图
	w4     W4Cmd   //四宫格
	w9     W9Cmd   //九宫格
}

type LiveCmd struct {
	src string //src 源rtsp
	dst string //dst 目标rtmp
}

type DbCmd struct {
	src   string //src 源rtsp
	dst   string //dst 目标rtmp
	start string //params
}

type JtCmd struct {
	Jtype string //计划截图，发起截图
	src   string //src 源rtsp
	dst   string //dst 目标rtmp
	start string //params
}

type W4Cmd struct {
	src1 string
	src2 string
	src3 string
	src4 string
	dst  string
}

type W9Cmd struct {
	src1 string
	src2 string
	src3 string
	src4 string
	src5 string
	src6 string
	src7 string
	src8 string
	src9 string
	dst  string
}

//func (l *Icmd) HandPlayRequest(conn *net.TCPConn) {
//	ipStr := conn.RemoteAddr().String()
//	defer func() {
//		fmt.Println("disconnected :" + ipStr)
//		conn.Close()
//	}()
//	buffer := make([]byte, 2048)
//	message, _ := conn.Read(buffer)
//	fmt.Println(string(buffer[:message]))

//	l = paser(buffer)
//	return
//	//	msg := time.Now().String() + "\n"
//	//	b := []byte(msg)
//	//	conn.Write(b)
//}

//func paser(cmd []byte) (lCmd *Icmd) {
//	if cmd[0] == 0x01 || cmd[0] == 0x03 {
//		return handLive(cmd)
//	}
//	if cmd[0] == 0x02 {
//		fmt.Println("点播")
//	}
//	if cmd[0] == 0x04 {
//		fmt.Println("保存")
//	}
//	return
//}

//func handLive(cmd []byte) (lCmd *Icmd) {
//	lCmd.op = cmd[0]
//	if cmd[1] == 0x01 {
//		lCmd.oc = 1
//	}
//	if cmd[1] == 0x02 {
//		lCmd.oc = 2
//	}
//	var ac uint8 = uint8(cmd[2])
//	var start uint8 = uint8(0x03) //
//	var stop uint8 = uint8(0x03)  //
//	stop = start + ac
//	var chanel []byte
//	for ; start < stop; start++ {
//		chanel = append(chanel, cmd[start])
//	}
//	//	fmt.Println("线路为: " + string(chanel))
//	lCmd.cqueue = string(chanel)
//	//name
//	ac = uint8(cmd[start])
//	start = start + 1
//	stop = start + ac
//	var name []byte
//	for ; start < stop; start++ {
//		name = append(name, cmd[start])
//	}
//	//	fmt.Println("名称为: " + string(name))
//	lCmd.name = string(name)
//	lCmd.dst = string(name)
//	//rtsp
//	ac = uint8(cmd[start])
//	start = start + 1
//	stop = start + ac
//	var rtsp []byte
//	for ; start < stop; start++ {
//		rtsp = append(rtsp, cmd[start])
//	}
//	//	fmt.Println("目标为: " + string(rtsp))
//	lCmd.src = string(rtsp)
//	return lCmd
//}
