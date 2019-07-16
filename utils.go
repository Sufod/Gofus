package main

import (
	"strconv"
	"strings"
)

func cryptPassword(password string, key string) string {
	chArray := [...]rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F',
		'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
		'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_'}
	str := []rune("#1")
	passwd := []rune(password)
	k := []rune(key)
	for i := 0; i < len(passwd); i++ {
		ch := passwd[i]
		ch2 := k[i]
		num2 := int(ch / '\u0010')
		num3 := int(ch % '\u0010')
		index := (num2 + int(ch2)) % len(chArray)
		num5 := (num3 + int(ch2)) % len(chArray)
		str = append(str, chArray[index], chArray[num5])
	}
	return string(str)
}

func (b Base64) encode64(data int) rune {
	return b.zkArray[data]
}

func (b Base64) decode64(data rune) int {
	return b.zkMap[data]
}

func ip2int(ip string) int {
	parts := strings.Split(ip, ".")
	ipAsInt := 0
	for i := 0; i < 4; i++ {
		tmp, _ := strconv.Atoi(parts[i])
		ipAsInt = ipAsInt | (tmp << uint(i*8))
	}
	return ipAsInt
}

func int2ip(ip int) string {
	return strconv.Itoa(ip>>24&0xff) + "." + strconv.Itoa(ip>>16&0xff) + "." + strconv.Itoa(ip>>8&0xff) + "." + strconv.Itoa(ip&0xff)
}

func (cypher DofusCypher) decodeIp(message string) string {
	obfIp := message[0:8]
	obfPort := message[8:11]
	ip := 0
	for i := 0; i < 8; i++ {
		ip = ip | ((int(obfIp[i]-48) & 15) << uint(4*(7-i)))
	}
	port := 0
	for i := 0; i < 3; i++ {
		port = port | ((cypher.base64.decode64(rune(obfPort[i])) & 63) << uint(6*(2-i)))
	}
	return int2ip(ip) + ":" + strconv.Itoa(port)
}

func (cypher DofusCypher) encodeIp(message string) string {
	parts := strings.Split(message, ":")
	ip := ip2int(parts[0])
	port, _ := strconv.Atoi(parts[1])

	obfIp := make([]rune, 8)
	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			obfIp[i+1] = rune(((ip >> uint(4*i)) & 15) + 48)
		} else {
			obfIp[i-1] = rune(((ip >> uint(4*i)) & 15) + 48)
		}
	}

	obfPort := make([]rune, 3)
	for i := 0; i < 3; i++ {
		obfPort[i] = cypher.base64.encode64((port >> uint(6*(2-i))) & 63)
	}

	return string(obfIp) + string(obfPort)
}

type Base64 struct {
	zkArray [64]rune
	zkMap   map[rune]int
}

type DofusCypher struct {
	base64 *Base64
}

func NewDofusCypher() *DofusCypher {
	c := &DofusCypher{
		base64: NewBase64(),
	}
	return c
}

func NewBase64() *Base64 {
	b := &Base64{
		zkMap: make(map[rune]int),
	}
	b.zkArray = [64]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_'}
	for index, r := range b.zkArray {
		b.zkMap[r] = index
	}
	return b
}

// private static int decode64(char data) {
// 	return zkArray.indexOf(data);
// }

// private static char encode64(int data) {
// 	return zkArray.get(data);
// }

// public static String decodeAXK(String rawPacket) {
// 	String obfIpPort = rawPacket.substring(3);
// 	String obfIp = obfIpPort.substring(0, 8);
// 	String obfPort = obfIpPort.substring(8, 11);
// 	int ip = 0;
// 	for (int i = 0; i < 8; i++) {
// 		int pos = 4 * (7 - i);
// 		ip |= (((obfIp.charAt(i) - 48) & 15) << pos);
// 	}
// 	int port = 0;
// 	for (int i = 0; i < 3; i++) {
// 		int pos = 6 * (2 - i);
// 		port |= (decode64(obfPort.charAt(i)) & 63) << pos;
// 	}
// 	return int2ip(ip) + ":" + port;
// }

// 	return new String(obfIp) + new String(obfPort);
// }

// private static String int2ip(int ip) {
// 	return IntStream.of(
// 			ip >> 24 & 0xff,
// 			ip >> 16 & 0xff,
// 			ip >> 8 & 0xff,
// 			ip & 0xff)
// 			.mapToObj(Integer::toString)
// 			.collect(Collectors.joining("."));
// }

// private static int ip2int(String ip) {
// 	String[] parts = ip.split("\\.");
// 	int iip = 0;
// 	for (int i = 0; i < parts.length; i++) {
// 		iip |= Integer.parseInt(parts[i]) << (i * 8);
// 	}
// 	return iip;
// }
