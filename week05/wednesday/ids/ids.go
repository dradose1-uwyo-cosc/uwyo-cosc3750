//Danny Radosevich
//COSC3750

//Examples on process/user/group IDS

package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func main() {
	//Get the user ID of the process
	uid := os.Getuid()
	fmt.Printf("User ID: %d\n", uid)
	//Get the group ID of the process
	gid := os.Getgid()
	fmt.Printf("Group ID: %d\n", gid)

	//get parent process ID
	ppid := os.Getppid()
	fmt.Printf("Parent Process ID: %d\n", ppid)
	//Get the process ID of the process
	pid := os.Getpid()
	fmt.Printf("Process ID: %d\n", pid)

	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		fmt.Printf("Error looking up user: %v\n", err)
		return
	}
	fmt.Printf("Username: %s\n", u.Username)

	g, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		fmt.Printf("Error looking up group: %v\n", err)
		return
	}
	fmt.Printf("Group Name: %s\n", g.Name)

}
