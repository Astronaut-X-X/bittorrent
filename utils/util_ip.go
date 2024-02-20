
func parseIp() {
	addr, err := net.ResolveUDPAddr("udp", node)
	if err != nil {
		fmt.Println(err.Error())
		continue
	}

	fmt.Println(addr.IP, addr.Port)
}