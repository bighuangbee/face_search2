package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"time"
)

func main() {

	var location, _ = time.LoadLocation("Asia/Shanghai")

	var stat unix.Stat_t
	err := unix.Stat("/hiar_face/registe_path/DSC08063.JPG", &stat)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印访问时间、修改时间和状态改变时间
	fmt.Println("Ctimespec", time.Unix(stat.Atim.Sec, stat.Atim.Nsec).String())
	fmt.Println("Mtim", time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).In(location).String())
	fmt.Println("Ctimespec", time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).String())

}
