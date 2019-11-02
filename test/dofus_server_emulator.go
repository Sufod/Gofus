package test

import (
	"net"
	"testing"
	"time"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/client"
	"github.com/Sufod/Gofus/internal/network"
	"gotest.tools/assert"
)

type DofusServerEmulator struct {
	*network.DofusSocket
}

func (emulator *DofusServerEmulator) Start(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:8081")              //Starting listening on local interface
	clientConn, _ := ln.Accept()                              //Blocking until a client connect
	emulator.DofusSocket = network.NewDofusSocket(clientConn) //Creating and initializing client socket conn
	go emulator.Listen()
	defer emulator.Close()
	emulator.handleClient(t)
}

func (emulator *DofusServerEmulator) handleClient(t *testing.T) {
	emulator.Send("HCzzybokxyrtkpjvxmmoxbnwiynojxdbqn")
	emulator.WaitForPacketAndAssertEqual(t, "1.30.0e")
	emulator.WaitForPacketAndAssertEqual(t, "testUser\n#1-haj_hNL00YR-9a75Y34YU3ZXX8f_6ZX")
	emulator.WaitForPacketAndAssertEqual(t, "Af")
	emulator.Send("Af2|3|0||-1")
	emulator.Send("AdtestUser")
	emulator.Send("Ac2")
	emulator.Send("AH601;1;75;1|605;1;75;1|609;1;75;1|604;1;75;1|608;1;75;1|603;1;75;1|607;1;75;1|611;1;75;1|602;1;75;1|606;1;75;1|610;1;75;1")
	emulator.Send("AlK0")
	emulator.Send("AQQuel+est+le+nom+de+mon+premier+animal+de+compagnie+%3F")
	emulator.WaitForPacketAndAssertEqual(t, "Ax")
	emulator.Send("AxK1609588265|608,2|608,2")
	emulator.WaitForPacketAndAssertEqual(t, "AX608")
	emulator.Send("AXK7?000001b-sa833b58")       // 7?000001b-s -> 127.0.0.1:8082
	ln, _ := net.Listen("tcp", "127.0.0.1:8082") //Starting listening on local interface
	clientConn, _ := ln.Accept()                 //Blocking until a client connect
	emulator.DofusSocket.Close()
	emulator.DofusSocket.Initialize(clientConn)
	go emulator.Listen()
	emulator.Send("HG")
	emulator.WaitForPacketAndAssertEqual(t, "ATa833b58")
	emulator.Send("ATK0")
	emulator.WaitForPacketAndAssertEqual(t, "Ak0")
	emulator.WaitForPacketAndAssertEqual(t, "AV")
	emulator.Send("BN")
	emulator.Send("AV0")
	emulator.WaitForPacketAndAssertEqual(t, "Agfr")
	emulator.WaitForPacketAndAssertEqual(t, "Ai3oMSVyiDZFcxq2zCBnGqdQ6")
	emulator.WaitForPacketAndAssertEqual(t, "AL")
	emulator.WaitForPacketAndAssertEqual(t, "Af")
	emulator.Send("BN")
	emulator.Send("BN")
	emulator.Send("BN")
	emulator.Send("Aq1")
	emulator.Send("ALK1609585713|2|80077248;Olibi;13;100;59ff7b;86fa7d;0;,1fe2,,,;0;608;;;|80087535;Testnini;1;71;-1;-1;-1;,,,,;0;608;;;")
	emulator.WaitForPacketAndAssertEqual(t, "AS80087535")
	emulator.WaitForPacketAndAssertEqual(t, "Af")
	emulator.Send("BN")
	emulator.Send("Af1|2|1|1|608")
	emulator.Send("Rx0")
	emulator.Send("ASK|80087535|Testnini|1|7|1|71|-1|-1|-1|45b5599~207~1~~;")
	emulator.Send("ZS0")
	emulator.Send("cC+*#%!$pi^")
	emulator.Send("al|270;0|49;1|319;0|98;0|147;0|466;0|245;0|515;2|294;0|73;1|122;1|171;1|441;0|220;0|490;0|269;0|48;1|318;0|97;0|146;0|465;0|244;0|514;1|23;2|293;0|72;2|121;2|170;0|440;0|219;0|268;0|47;1|317;0|96;0|145;0|464;0|243;0|513;1|22;0|292;0|71;0|120;1|169;1|218;0|488;0|267;0|46;1|316;0|95;0|144;0|463;0|512;1|21;0|291;0|70;0|119;1|168;0|217;0|487;0|266;0|536;0|45;1|315;0|94;0|143;2|462;0|511;2|20;0|290;0|69;1|339;0|118;2|167;0|216;0|486;0|44;1|314;0|93;0|461;0|510;2|19;0|289;0|68;1|338;0|117;1|166;0|215;0|485;0|43;1|313;0|92;0|141;0|460;0|509;2|18;0|288;0|67;1|337;0|116;2|165;0|214;0|484;0|42;0|312;0|91;0|140;0|459;0|508;2|17;0|287;0|66;1|336;0|115;1|164;0|213;0|483;0|41;0|311;0|139;0|507;2|16;0|286;0|65;1|335;0|114;2|163;0|212;0|482;1|261;0|40;0|310;0|89;0|138;0|457;0|236;0|506;2|15;0|285;0|64;1|334;1|113;1|162;0|211;0|481;2|260;0|39;0|309;0|88;0|137;0|235;1|505;2|14;0|284;0|63;1|333;0|112;2|161;0|210;0|480;1|259;0|38;1|308;0|87;0|136;0|455;0|234;0|504;2|13;0|62;1|332;0|111;1|209;0|479;0|258;0|37;1|307;0|86;0|135;0|454;0|233;1|503;2|12;2|61;2|331;0|110;0|159;0|208;0|478;1|257;0|306;0|85;0|134;0|453;0|232;0|502;2|11;2|281;0|60;0|330;0|109;0|158;0|207;0|477;1|256;0|35;0|84;0|133;0|182;0|452;0|231;0|501;0|10;1|280;2|59;1|329;0|108;0|157;0|206;0|476;2|255;0|34;0|304;0|83;0|132;0|181;0|451;0|230;0|500;0|9;0|279;1|328;0|107;1|156;0|205;0|254;0|33;0|303;0|82;0|131;0|180;0|450;0|229;0|499;0|8;0|278;0|57;0|327;0|106;1|155;0|204;0|474;0|253;1|32;2|302;0|81;1|130;0|179;0|449;0|228;0|498;0|7;0|277;0|56;0|326;0|105;0|154;0|203;0|473;0|252;0|31;1|301;0|80;1|129;0|178;0|448;0|227;0|497;0|6;0|276;0|55;0|325;0|153;0|202;0|472;0|251;0|30;0|300;0|79;1|128;0|177;0|447;0|226;0|496;0|5;0|275;0|54;1|324;0|103;2|152;0|201;0|471;0|250;0|29;0|299;0|78;0|127;0|446;0|225;0|495;0|4;0|274;0|53;2|323;0|102;0|151;0|200;0|470;0|249;0|28;0|298;0|77;0|126;0|175;0|445;0|224;0|494;0|3;0|273;0|322;0|101;0|150;0|469;0|248;0|27;2|297;0|76;1|125;0|174;0|444;0|223;0|493;0|2;0|272;0|51;1|321;0|100;0|149;0|468;1|247;0|26;0|296;0|75;2|124;0|173;0|443;0|222;0|492;2|1;0|271;0|50;1|320;0|99;0|148;0|467;0|246;0|25;0|295;0|74;1|123;2|442;0|221;0|491;0|0;0")
	emulator.Send("SLo+")
	emulator.Send("SL121~1~b;125~1~c;128~1~d;")
	emulator.Send("AR6bk")
	emulator.Send("Ow1|1000")
	emulator.Send("FO-")
	emulator.Send("Im189")
	emulator.Send("Im0152;2019~11~01~19~56~2.7.208.108")
	emulator.Send("Im0153;93.22.150.248")
	emulator.Send("Im01;6")
	emulator.WaitForPacketAndAssertEqual(t, "GC1")

}

//WaitForPacket blocks until a message is available to read in the channel
func (emulator *DofusServerEmulator) WaitForPacketAndAssertEqual(t *testing.T, expectedPacket string) {
	packet, err := emulator.WaitForPacket()
	assert.NilError(t, err, "Expected paquet "+expectedPacket)
	receivedPacket := packet[:len(packet)-2]
	assert.Equal(t, receivedPacket, expectedPacket)
}

func (emulator *DofusServerEmulator) startClient() {
	var cfg configs.ConfigHolder = configs.ConfigHolder{
		DofusAuthServer: "127.0.0.1:8081",
		DofusServerName: "Ayuto",
		DofusVersion:    "1.30.0e",
		Credentials: &configs.Credentials{
			Username: "testUser",
			Password: "MonSUperp4ssword",
		},
	}

	time.Sleep(1 * time.Second)
	client := client.NewDofusClient(cfg)
	client.Start()
}
