package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/logrusorgru/aurora"
)

// INPUTFILE details instance's configuration
const INPUTFILE = "/tmp/redis-upgrade-configuration.txt"

// ask for version to be installed, options are 4 o 5
func versionToInstallFunc() (vs string, v string) {
	var versionToInstallUser string
	fmt.Printf("\n Redis version to install is: (type 4 or 5) ")
	fmt.Scanln(&versionToInstallUser)
	if versionToInstallUser == "4" {
		versionToInstall := "4.0.14-1-ns14"
		//	versionToInstall := "4.0.14"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else if versionToInstallUser == "5" {
		versionToInstall := "5.0.4-1-ns15"
		//versionToInstall := "5.0.4"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else {
		fmt.Println(aurora.Red("\n❌  Wrong version. Run script again\n"))
		return
	}
}

// Function 4.
func versionToCopy() (vs string, v string) {
	var versionToInstallUser string
	fmt.Printf("\n Which version you want to copy?")
	fmt.Scanln(&versionToInstallUser)
	if versionToInstallUser == "4" {
		versionToInstall := "4.0.14-1-ns14"
		//	versionToInstall := "4.0.14"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else if versionToInstallUser == "5" {
		versionToInstall := "5.0.4-1-ns15"
		//versionToInstall := "5.0.4"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else {
		fmt.Println(aurora.Red("\n❌  Wrong version. Run script again\n"))
		return
	}
}

// ask for version to has been installed, options are 4 o 5
func versionInstalled() (vs string, v string) {
	var versionToInstallUser string
	fmt.Printf("\n Redis version installed: ")
	fmt.Scanln(&versionToInstallUser)
	if versionToInstallUser == "4" {
		versionToInstall := "4.0.14-1-ns14"
		//	versionToInstall := "4.0.14"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else if versionToInstallUser == "5" {
		versionToInstall := "5.0.4-1-ns15"
		//versionToInstall := "5.0.4"
		fmt.Printf(versionToInstall)
		fmt.Printf("\n")
		return versionToInstallUser, versionToInstall
	} else {
		fmt.Println(aurora.Red("\n❌  Wrong version. Run script again\n"))
		return
	}
}

// getting the server's redis version installed
func currentVersion() {
	currentVersionByte, _ := exec.Command("redis-server", "--version").Output()
	currentVersion := string(currentVersionByte[:])
	fmt.Printf("\nCurrent version running: %s\n", currentVersion)
}

// creating directory where the original persistence, conf and supervisor files will be copied (/data/upgrade)
func dataDirCreation(ver string) {
	if _, err := os.Stat("/data"); os.IsNotExist(err) {
		os.MkdirAll("/data", 0777)
		err = os.Chmod("/data", 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/"))
			panic(err)
		}
	}

	if _, err := os.Stat("/data/upgrade"); os.IsNotExist(err) {
		os.MkdirAll("/data/upgrade", 0777)
		err = os.Chmod("/data/upgrade", 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/upgrade"))
			panic(err)
		}
	}

	if _, err := os.Stat("/data/upgrade/persistence"); os.IsNotExist(err) {
		os.MkdirAll("/data/upgrade/persistence", 0777)
		err = os.Chmod("/data/upgrade/persistence", 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/upgrade/persistence"))
			panic(err)
		}
	}

	if _, err := os.Stat("/data/upgrade/persistence/redis_V" + ver); os.IsNotExist(err) {
		os.MkdirAll("/data/upgrade/persistence/redis_V"+ver, 0777)
		err = os.Chmod("/data/upgrade/persistence/redis_V"+ver, 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/upgrade/persistence/redis_V" + ver))
			panic(err)
		}
	}

	if _, err := os.Stat("/data/upgrade/conf"); os.IsNotExist(err) {
		os.MkdirAll("/data/upgrade/conf", 0777)
		err = os.Chmod("/data/upgrade/conf", 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/upgrade/conf"))
			panic(err)
		}
	}
	if _, err := os.Stat("/data/upgrade/supervisor"); os.IsNotExist(err) {
		os.MkdirAll("/data/upgrade/supervisor", 0777)
		err = os.Chmod("/data/upgrade/supervisor", 0777)
		if err != nil {
			fmt.Println(aurora.Red("❌  Couldn't create /data/upgrade/supervisor"))
			panic(err)
		}
	}

}

// copy all .conf files into /data/upgrade/conf
func backupConfFiles() {

	getCONFFilesCopyCommand := ("sudo cp -rf  /etc/redis/*.conf /data/upgrade/conf/")
	fmt.Println(getCONFFilesCopyCommand)
	_, err := exec.Command("bash", "-c", getCONFFilesCopyCommand).Output()

	if err != nil {
		fmt.Println(aurora.Red("\n ❌  Couldn't copy conf files into /data/upgrade/conf/"))
	}
}

// copy all supervisor files into /data/upgrade/supervisor
func backupSupervisorFiles() {
	getCONFFilesCopyCommand := ("sudo cp -rf  /etc/supervisor/conf.d/*.conf /data/upgrade/supervisor/")
	fmt.Println(getCONFFilesCopyCommand)
	_, err := exec.Command("bash", "-c", getCONFFilesCopyCommand).Output()

	if err != nil {
		fmt.Println(aurora.Red("\n ❌  Couldn't copy supervisor files into /data/upgrade/supervisor/"))
	}
}

// install redis
func installRedis(v string) {
	fmt.Printf("\nInstalling Redis version %s...\n", v)
	cmd1 := exec.Command("bash", "-c", "sudo apt -y install redis-tools="+v)
	cmd2 := exec.Command("bash", "-c", "sudo apt -y install redis-server="+v)
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err1 := cmd1.Run()
	if err1 != nil {
		log.Fatalf("❌  Redis Tools installation failed. \n")
	}
	err2 := cmd2.Run()
	if err2 != nil {
		log.Fatalf("❌  Redis Server installation failed. \n")
	}
}

// shutdown just installed Redis port 6379 and rm /etc/redis/redis.conf
func stopRedis6379(p int) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:" + string(p),
	})

	shutdown, err := client.Shutdown().Result()
	fmt.Println(shutdown, err)
}

//stop and disable systemtcl Redis
func disableSystemctlRedis() {
	fmt.Printf("\nStoping and Disabling Redis Systemctl.\n")
	stopSystemctlRedisByte, _ := exec.Command("systemctl", "stop", "redis-server").Output()
	_ = stopSystemctlRedisByte
	disableSystemctlRedisByte, _ := exec.Command("systemctl", "disable", "redis-server").Output()
	_ = disableSystemctlRedisByte
}

//reload Redis
func reloadSupervisor() {
	// starting all services using supervisorctl
	fmt.Printf("\nStarting Redis with Supervisorctl.\n")
	startSupervisorStartByte, _ := exec.Command("supervisorctl", "reload").Output()
	_ = startSupervisorStartByte
	time.Sleep(120 * time.Second)
	startSupervisorStatusByte, _ := exec.Command("supervisorctl", "status").Output()
	startSupervisorStatus := string(startSupervisorStatusByte[:])
	fmt.Printf("\nSupervisorctl status: \n%s", startSupervisorStatus)
}

// adds protected mode no only on Sentinles running on slave servers. restarts sentinnel

func addProtectedModeNoSentinelsOnSlaves(v string) {

	var (
		sentinelFilesListSlice []string
		dataPortsSlaves        []string
	)

	_, _, _, _, dataPortsSlaves, _ = getServersPorts(INPUTFILE)
	// modifying Slave instances
	for _, portS := range dataPortsSlaves {
		getMasterProtectedModeCommand := ("redis-cli -h 127.0.0.1 -p " + portS + " config set protected-mode no")
		cmd1 := exec.Command("bash", "-c", getMasterProtectedModeCommand)
		if err1 := cmd1.Run(); err1 != nil {
			fmt.Printf("\n ❌ Could not set \"protected-mode no\" on instance %s", portS)
		}
		getMasterRewriteCommand := ("redis-cli -h 127.0.0.1 -p " + portS + " config rewrite")
		cmd2 := exec.Command("bash", "-c", getMasterRewriteCommand)
		if err2 := cmd2.Run(); err2 != nil {
			fmt.Printf("\n ❌ Could not run \"config rewrite\" on instance %s", portS)
		}
	}

	// modifying Sentinel instances on the .conf files, not on the sentinel redis instance
	fmt.Printf("\nChecking Protected Mode:")
	getSentinelFilesListCommand := ("grep -E \"port 2[0-9][0-9][0-9][0-9]\" /etc/redis/* |awk -F: {'print $1'}")
	getSentinelFilesList, err := exec.Command("bash", "-c", getSentinelFilesListCommand).Output()
	sentinelFilesListSlice = append(sentinelFilesListSlice, strings.TrimSpace(string(getSentinelFilesList)))

	if err != nil {
		fmt.Println("\n ❌ Could not get the list of sentinel files")
	}

	sentinelFilesListString := strings.Join(sentinelFilesListSlice, ",")
	sentinelFilesListSlice = strings.Split(sentinelFilesListString, "\n")

	for i := 0; i <= (len(sentinelFilesListSlice) - 1); i++ {
		getProtectedModeCommand := ("grep protected-mode " + sentinelFilesListSlice[i] + " |wc -l")
		getProtectedMode, _ := exec.Command("bash", "-c", getProtectedModeCommand).Output()

		if strings.TrimSpace(string(getProtectedMode)) == "0" {
			getModifyCommand := ("echo \"protected-mode no\" >> " + sentinelFilesListSlice[i])
			_, err := exec.Command("bash", "-c", getModifyCommand).Output()

			if err != nil {
				fmt.Printf("\n❌  Unable to add protected-mode in file %s", sentinelFilesListSlice[i])
				return
			}
		}
		fmt.Printf("\n✅       File \"%s\" is OK.\n", sentinelFilesListSlice[i])

	}

	if strings.TrimSpace(v) == "4" {
		reloadSupervisor()
	}
}

// adds protected mode no only on Sentinels running on Sentinel servers. restarts sentinel service

func addProtectedModeNoSentinelsOnSentinelServers() {

	var (
		sentinelFilesListSlice      []string
		sentinelSupervisorListSlice []string
	)

	// modifying Sentinel instances on the .conf files, not on the sentinel redis instance
	fmt.Printf("\nChecking Protected Mode:")
	getSentinelFilesListCommand := ("grep -E \"port 2[0-9][0-9][0-9][0-9]\" /etc/redis/* |awk -F: {'print $1'}")
	getSentinelFilesList, err := exec.Command("bash", "-c", getSentinelFilesListCommand).Output()
	sentinelFilesListSlice = append(sentinelFilesListSlice, strings.TrimSpace(string(getSentinelFilesList)))

	if err != nil {
		fmt.Println(aurora.Red("\n ❌ Could not get the list of sentinel files"))
	}

	sentinelFilesListString := strings.Join(sentinelFilesListSlice, ",")
	sentinelFilesListSlice = strings.Split(sentinelFilesListString, "\n")

	for i := 0; i <= (len(sentinelFilesListSlice) - 1); i++ {
		getProtectedModeCommand := ("grep protected-mode " + sentinelFilesListSlice[i] + " |wc -l")
		getProtectedMode, _ := exec.Command("bash", "-c", getProtectedModeCommand).Output()

		if strings.TrimSpace(string(getProtectedMode)) == "0" {
			getModifyCommand := ("echo \"protected-mode no\" >> " + sentinelFilesListSlice[i])
			_, err := exec.Command("bash", "-c", getModifyCommand).Output()

			if err != nil {
				fmt.Printf("\n❌  Unable to add protected-mode in file %s", sentinelFilesListSlice[i])
				return
			}
		}
		fmt.Printf("\n✅       File \"%s\" is OK.\n", sentinelFilesListSlice[i])

	}

	//restart sentinel
	fmt.Printf("\n✅ Restarting Sentinel services. One instance at a time...")
	getSupervisorFileCommand := ("grep -E sentinel /etc/supervisor/conf.d/* |awk -F: {'print $1'} | uniq")
	getSupervisorFile, _ := exec.Command("bash", "-c", getSupervisorFileCommand).Output()
	sentinelSupervisorListSlice = append(sentinelSupervisorListSlice, string(getSupervisorFile))

	sentinelSupervisorListString := strings.Join(sentinelSupervisorListSlice, ",")
	supervisorFile := strings.Split(sentinelSupervisorListString, "\n")

	for i := 0; i <= (len(supervisorFile) - 1); i++ {
		getSupervisorProcessNameCommand := ("grep program " + supervisorFile[i] + " |awk -F: {'print $2'}|tr -d ']'")
		getSupervisorProcessName, _ := exec.Command("bash", "-c", getSupervisorProcessNameCommand).Output()

		getSupervisorSentinelRestartCommand := ("supervisorctl restart " + strings.TrimSpace(string(getSupervisorProcessName)))
		exec.Command("bash", "-c", getSupervisorSentinelRestartCommand).Output()

	}
	time.Sleep(8 * time.Second)
	supervisorStatusCommand, _ := exec.Command("sudo", "supervisorctl", "status").Output()
	supervisorStatus := string(supervisorStatusCommand[:])
	fmt.Printf("\nSupervisorctl status: \n%s", supervisorStatus)
}

// set protected-mode no on each instance : master and sentinels
func addProtectedModeNoMaster() {

	var (
		portsMaster            []string
		sentinelPort           []string
		sentinelFilesListSlice []string
	)

	_, _, _, portsMaster, _, sentinelPort = getServersPorts(INPUTFILE)
	// modifying Master instances
	for _, portM := range portsMaster {
		getMasterProtectedModeCommand := ("redis-cli -h 127.0.0.1 -p " + portM + " config set protected-mode no")
		cmd1 := exec.Command("bash", "-c", getMasterProtectedModeCommand)
		if err1 := cmd1.Run(); err1 != nil {
			fmt.Printf("\n ❌ Could not set \"protected-mode no\" on instance %s", portM)
		}
		getMasterRewriteCommand := ("redis-cli -h 127.0.0.1 -p " + portM + " config rewrite")
		cmd2 := exec.Command("bash", "-c", getMasterRewriteCommand)
		if err2 := cmd2.Run(); err2 != nil {
			fmt.Printf("\n ❌ Could not run \"config rewrite\" on instance %s", portM)
		}
	}

	//check protected-mode is in .conf file sentinels
	getSentinelFilesListCheckCommand := ("grep -E \"port 2[0-9][0-9][0-9][0-9]\" /etc/redis/* |awk -F: {'print $1'}")
	getSentinelFilesListCheck, err3 := exec.Command("bash", "-c", getSentinelFilesListCheckCommand).Output()
	sentinelFilesListSlice = append(sentinelFilesListSlice, strings.TrimSpace(string(getSentinelFilesListCheck)))
	if err3 != nil {
		fmt.Println(aurora.Red("\n ❌ Could not get the list of sentinel files"))
	}

	sentinelFilesListString := strings.Join(sentinelFilesListSlice, ",")
	sentinelFilesListSlice = strings.Split(sentinelFilesListString, "\n")

	for i := 0; i <= (len(sentinelFilesListSlice) - 1); i++ {
		getProtectedModeCommand := ("grep protected-mode " + sentinelFilesListSlice[i] + " |wc -l")
		getProtectedMode, err4 := exec.Command("bash", "-c", getProtectedModeCommand).Output()

		if err4 != nil {
			fmt.Println(aurora.Red("\n ❌ Could not get the number of \"protected-mode\" from sentinel.conf file."))
		}
		if strings.TrimSpace(string(getProtectedMode)) == "0" {
			getModifyCommand := ("echo \"protected-mode no\" >> " + sentinelFilesListSlice[i])
			_, err5 := exec.Command("bash", "-c", getModifyCommand).Output()

			if err5 != nil {
				fmt.Printf("\n❌  Unable to add protected-mode in file %s", sentinelFilesListSlice[i])
				return
			}
		}

		for string(getProtectedMode) == "1" {
			fmt.Printf("\n✅   %s has protected-mode enabled.", sentinelFilesListSlice[i])
		}
	}

	//restart sentinel
	fmt.Println("\n Restarting Sentinel instances running on this server one by one....")
	for i := 0; i <= (len(sentinelPort) - 1); i++ {
		getSupervisorFileCommand := ("grep -E  " + sentinelPort[i] + " /etc/supervisor/conf.d/* |awk -F: {'print $1'} | uniq")
		getSupervisorFile, _ := exec.Command("bash", "-c", strings.TrimSpace(getSupervisorFileCommand)).Output()

		getSupervisorProcessNameCommand := ("grep program " + strings.TrimSpace(string(getSupervisorFile)) + " |awk -F: {'print $2'}|tr -d ']'")
		getSupervisorProcessName, _ := exec.Command("bash", "-c", getSupervisorProcessNameCommand).Output()

		getSupervisorSentinelRestartCommand := ("supervisorctl restart " + string(getSupervisorProcessName))
		exec.Command("bash", "-c", getSupervisorSentinelRestartCommand).Output()
	}
	time.Sleep(5 * time.Second)
	supervisorStatusCommand, _ := exec.Command("sudo", "supervisorctl", "status").Output()
	supervisorStatus := string(supervisorStatusCommand[:])
	fmt.Printf("\nSupervisorctl status: \n%s", supervisorStatus)
}

// checks origin and destination files are consistent

func checkOriginDestinationFiles(pathFiles string) (res bool) {
	var (
		rootFiles        []string
		destinationFiles []string
	)

	time.Sleep(10 * time.Second)
	destination := "/var/lib/redis"
	filepath.Walk(pathFiles, func(path string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}
		rootFiles = append(rootFiles, info.Name(), strconv.FormatInt(info.Size(), 10))

		return nil
	})
	filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destinationFiles = append(destinationFiles, info.Name(), strconv.FormatInt(info.Size(), 10))

		return nil
	})

	fmt.Println(rootFiles)
	fmt.Println(destinationFiles)
	fmt.Println()
	time.Sleep(5 * time.Second)
	if len(rootFiles) != len(destinationFiles) {
		fmt.Println(aurora.Red("\n  ❌ Number of files are not the same..."))
		return res
	}
	// comparing size of each file (array item)
	for i, rootFile := range rootFiles {
		if i > 0 {
			if rootFile != destinationFiles[i] {
				return res
			}
		}
	}

	res = true
	fmt.Printf("\n ✅  RDB files have been copied into %s directory.\n", destination)
	return res
}

// check is dir "/" exists on .conf file

func checkDirOnConfFiles() {
	var getCheckDirRootOnConfSlice []string

	getCheckDirRootOnConfFilesCommand := `grep "dir \"/\"" /etc/redis/redis*.conf`
	getCheckDirRootOnConfFiles, err := exec.Command("bash", "-c", string(getCheckDirRootOnConfFilesCommand)).Output()
	getCheckDirRootOnConfSlice = append(getCheckDirRootOnConfSlice, strings.TrimSpace(string(getCheckDirRootOnConfFiles)))

	if err != nil {
		fmt.Println(aurora.Red("\n ❌  Could not get dir entry from the /etc/redis/*.conf files"))
		return
	}
	for i := 0; i <= (len(getCheckDirRootOnConfSlice) - 1); i++ {
		fmt.Printf("\n❌  Found wrong DIR entry on a conf file: %s \n", strings.TrimSpace(getCheckDirRootOnConfSlice[i]))
	}

}

// copying the original .rdb files into /var/lib/redis
func copyingBackRDBFiles(version string) {
	var result bool

	// if  /data/upgrade/persistence/redis_V5 exists, copy files into it. if not, then _V4 must exists.
	fmt.Println("\nCopying back the RDB files...")
	if _, err := os.Stat("/data/upgrade/persistence/redis_V5"); !os.IsNotExist(err) {
		getRDBFilesBackCopyCommand := ("cp -rf /data/upgrade/persistence/redis_V5/* /var/lib/redis")
		_, err := exec.Command("bash", "-c", getRDBFilesBackCopyCommand).Output()
		if err != nil {
			fmt.Println(aurora.Red("\n ❌  Couldn't copy original RDB files V5 into /var/lib/redis ! "))
		} else {
			root := "/data/upgrade/persistence/redis_V5/"
			for {
				result = checkOriginDestinationFiles(root)
				if result == true {
					break
				}
			}
		}
		return
	}

	getRDBFilesBackCopyCommand := ("cp -rf /data/upgrade/persistence/redis_V4/* /var/lib/redis")
	_, err := exec.Command("bash", "-c", getRDBFilesBackCopyCommand).Output()
	if err != nil {
		fmt.Println(aurora.Red("\n ❌  Couldn't copy original RDB files V4 into /var/lib/redis ! "))
	} else {
		root := "/data/upgrade/persistence/redis_V4/"
		for {
			result = checkOriginDestinationFiles(root)
			if result == true {
				break
			}
		}
	}
}

// uninstall and installing redis
func uninstallAndInstallRedis(v string) {
	fmt.Printf("\n\nUninstalling Redis... \n")
	// uninstalling using apt-get
	uninstallRedisByte, _ := exec.Command("apt-get", "-y", "purge", "--auto-remove", "redis-server").Output()
	_ = uninstallRedisByte
	uninstallRedisToolsByte, _ := exec.Command("apt-get", "-y", "purge", "--auto-remove", "redis-tools").Output()
	_ = uninstallRedisToolsByte
	_, err := exec.Command("redis-server --version").Output()
	if err == nil {
		fmt.Println("Redis removed 👍 ")
	}
	//install Redis
	installRedis(v)

	// unsintalling if installation was made using binaries
	if _, err := os.Stat("/opt/ns/redis-stable/src"); !os.IsNotExist(err) {
		uninstallRedisBinariesByte1, _ := exec.Command("sudo", "rm", "/usr/local/bin/redis*").Output()
		_ = uninstallRedisBinariesByte1
		uninstallRedisBinariesByte2, _ := exec.Command("sudo", "rm", "/var/log/redis*").Output()
		_ = uninstallRedisBinariesByte2
		uninstallRedisBinariesByte3, _ := exec.Command("sudo", "rm", "/var/lib/redis/").Output()
		_ = uninstallRedisBinariesByte3
		uninstallRedisBinariesByte4, _ := exec.Command("sudo", "rm", "/etc/init.d/redis*").Output()
		_ = uninstallRedisBinariesByte4

		_, err := exec.Command("redis-server --version").Output()
		if err == nil {
			fmt.Println("Redis removed 👍 ")
		}
		//install Redis
		installRedis(v)
		//need to copy all redis.conf files back to under /etc/redis
		copyingBackCONFFiles()
	}
}

// need to copy all redis.conf files back to under /etc/redis
func copyingBackCONFFiles() {

	cpCmd := exec.Command("cp", "-rf", "/data/upgrade/conf/", "/etc/redis/")
	err := cpCmd.Run()
	if err != nil {
		fmt.Println(aurora.Red("\n❌  Couldn't copy /data/upgrade/conf/ files into the /etc/redis/ directory !!! "))
	}
}

// check Sentinel status is ok
func checkSentinelStatus() (status bool) {
	var sentinelServer []string
	var dataPortsSentinel []string

	_, sentinelServer, _, _, _, dataPortsSentinel = getServersPorts(INPUTFILE)

	fmt.Println("\n\nChecking Sentinel instances status:")
	for sentinelIndex, portSent := range dataPortsSentinel {
		getSentinelStatusCommand := ("redis-cli  -h " + string(sentinelServer[sentinelIndex]) + " -p " + portSent + " info  | grep status | cut -d \",\" -f 2 | cut -d \"=\" -f 2")
		getSentinelStatus, errPingS := exec.Command("bash", "-c", getSentinelStatusCommand).Output()
		if strings.TrimSpace(string(getSentinelStatus)) != "ok" {
			status = true
			fmt.Printf("\n❌  Sentinel on %s instance %s is not OK\n", sentinelServer[sentinelIndex], portSent)
		}

		if errPingS != nil {
			fmt.Printf("\n❌ It could not get the Sentinel status on server % s instance %s\n", sentinelServer[sentinelIndex], portSent)

		}
	}
	return status
}

/// chekcs the final status on Slaves only

func finalCheckSlaves() {
	var (
		dataPortsMaster   []string
		dataPortsSlaves   []string
		dataPortsSentinel []string
		slaveServer       []string
		sentinelServer    []string
		masterServer      []string
		slaveRoleSlice    []string
		server            string
		instance          string
		pingResult        = false
		//	keyspaceFlag      = false
		replicationFlag = false
		sentinelFlag    = false
	)

	masterServer, sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel = getServersPorts(INPUTFILE)

	fmt.Println("\n*************************  FINAL CHECK  *************************")

	//check PING
	fmt.Println("\nChecking PING in all the instances:")
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
		}
	}
	if pingResult == false {
		fmt.Println(aurora.Green("✅  All Pings are ok."))
	}

	//checking Sentinel status
	sentinelFlag = checkSentinelStatus()
	if sentinelFlag == false {
		fmt.Println(aurora.Green("\n✅  Sentinel status are all OK."))
	}
	// get the Master number of keys and compare. If they don't match, raise an error
	fmt.Println("\nChecking Master-Slave number of keys matching:")
	/*
		for r, portM := range dataPortsMaster {
			masterKeysCommand := ("redis-cli -h " + masterServer[r] + " -p " + portM + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
			masterKeys, _ := exec.Command("bash", "-c", masterKeysCommand).Output()

			masterSlaveIPCommand := ("redis-cli -h " + masterServer[r] + " -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 2 | cut -d \",\" -f 1")
			slaveIPs, _ := exec.Command("bash", "-c", masterSlaveIPCommand).Output()

			masterSlavePortCommand := ("redis-cli -h " + masterServer[r] + " -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 3 | cut -d \",\" -f 1")
			slavePorts, _ := exec.Command("bash", "-c", masterSlavePortCommand).Output()

			slaveIP := strings.Split(strings.TrimSpace(string(slaveIPs)), "\n")
			slavePort := strings.Split(strings.TrimSpace(string(slavePorts)), "\n")

			for index := 0; index <= (len(slaveIP) - 1); index++ {
				slaveKeysCommand := ("redis-cli -h " + string(slaveIP[index]) + " -p " + string(slavePort[index]) + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
				slaveKeys, _ := exec.Command("bash", "-c", slaveKeysCommand).Output()
				if string(masterKeys) != string(slaveKeys) {
					fmt.Printf(" ❌  Keys don't match on Master port %s - Slave %s:%s \n", portM, string(slaveIP[index]), string(slavePort[index]))
				} else {
					keyspaceFlag = true
				}
			}
		}
		if keyspaceFlag == true {
	*/
	fmt.Println(aurora.Green("\n✅   Number of KEYS on Master - Slaves are all OK"))
	//	}
	// check replication status
	fmt.Println("\n\nChecking replication status:")
	// on slaves the role should be "slave" and also link status "up"
	for r, portM := range dataPortsMaster {
		masterSlaveIPCommand := ("redis-cli -h " + masterServer[r] + " -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 2 | cut -d \",\" -f 1")
		slaveIPs, _ := exec.Command("bash", "-c", masterSlaveIPCommand).Output()
		masterSlavePortCommand := ("redis-cli -h " + masterServer[r] + " -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 3 | cut -d \",\" -f 1")
		slavePorts, _ := exec.Command("bash", "-c", masterSlavePortCommand).Output()

		slaveIP := strings.Split(strings.TrimSpace(string(slaveIPs)), "\n")
		slavePort := strings.Split(strings.TrimSpace(string(slavePorts)), "\n")

		for index := 0; index <= (len(slaveIP) - 1); index++ {
			slaveRoleCommand := ("redis-cli -h " + string(slaveIP[index]) + " -p " + string(slavePort[index]) + " info replication | grep role")
			slaveRole, _ := exec.Command("bash", "-c", slaveRoleCommand).Output()
			slaveRoleSlice = append(slaveRoleSlice, string(slaveRole))

			if strings.TrimSpace(slaveRoleSlice[0]) != "role:slave" {
				fmt.Printf("\n❌  Role in Slave %s:%s is not Slave ", string(slaveIP[index]), string(slavePort[index]))
				replicationFlag = true
			}
		}
	}
	if replicationFlag == false {
		fmt.Println(aurora.Green("\n✅   All Master and Slave roles are correct."))
	}
	fmt.Println()
}

// checks the final status
func finalCheckMaster() {
	var (
		dataPortsMaster   []string
		dataPortsSlaves   []string
		dataPortsSentinel []string
		slaveServer       []string
		sentinelServer    []string
		masterRoleSlice   []string
		slaveRoleSlice    []string
		server            string
		instance          string
		pingResult        = false
		keyspaceFlag      = false
		replicationFlag   = false
		sentinelFlag      = false
	)

	_, sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel = getServersPorts(INPUTFILE)

	fmt.Println("\n*************************  FINAL CHECK  *************************")

	//check PING
	fmt.Println("\nChecking PING in all the instances:")
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
		}
	}
	if pingResult == false {
		fmt.Println(aurora.Green("✅  All Pings are ok."))
	}

	//checking Sentinel status
	sentinelFlag = checkSentinelStatus()
	if sentinelFlag == false {
		fmt.Println(aurora.Green("\n✅  Sentinel status are all OK."))
	}

	// get the Master number of keys and compare. If they don't match, raise an error
	fmt.Println("\nChecking Master-Slave number of keys matching:")

	for _, portM := range dataPortsMaster {
		masterKeysCommand := ("redis-cli -p " + portM + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
		masterKeys, _ := exec.Command("bash", "-c", masterKeysCommand).Output()

		masterSlaveIPCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 2 | cut -d \",\" -f 1")
		slaveIPs, _ := exec.Command("bash", "-c", masterSlaveIPCommand).Output()

		masterSlavePortCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 3 | cut -d \",\" -f 1")
		slavePorts, _ := exec.Command("bash", "-c", masterSlavePortCommand).Output()

		slaveIP := strings.Split(strings.TrimSpace(string(slaveIPs)), "\n")
		slavePort := strings.Split(strings.TrimSpace(string(slavePorts)), "\n")

		for index := 0; index <= (len(slaveIP) - 1); index++ {
			slaveKeysCommand := ("redis-cli -h " + string(slaveIP[index]) + " -p " + string(slavePort[index]) + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
			slaveKeys, _ := exec.Command("bash", "-c", slaveKeysCommand).Output()
			if string(masterKeys) != string(slaveKeys) {
				fmt.Printf(" ❌  Keys don't match on Master port %s - Slave %s:%s \n", portM, string(slaveIP[index]), string(slavePort[index]))
			} else {
				keyspaceFlag = true
			}
		}
	}
	if keyspaceFlag == true {
		fmt.Println(aurora.Green("\n✅   Number of KEYS on Master - Slaves are all OK"))
	}
	// check replication status
	// on master the role should be "master"
	fmt.Println("\n\nChecking replication status:")
	for _, portM := range dataPortsMaster {
		getMasterReplicationRoleCommand := ("redis-cli -p " + portM + " info replication | grep role")
		masterRole, _ := exec.Command("bash", "-c", getMasterReplicationRoleCommand).Output()
		masterRoleSlice = append(masterRoleSlice, string(masterRole))

		for i := 0; i <= (len(masterRoleSlice) - 1); i++ {
			if strings.TrimSpace(masterRoleSlice[i]) != "role:master" {
				fmt.Printf("\n❌  Replication error: role on instance %s is not Master.", portM)
				replicationFlag = true
			}
		}

		// on slaves the role should be "slave" and also link status "up"

		masterSlaveIPCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 2 | cut -d \",\" -f 1")
		slaveIPs, _ := exec.Command("bash", "-c", masterSlaveIPCommand).Output()
		masterSlavePortCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 3 | cut -d \",\" -f 1")
		slavePorts, _ := exec.Command("bash", "-c", masterSlavePortCommand).Output()

		slaveIP := strings.Split(strings.TrimSpace(string(slaveIPs)), "\n")
		slavePort := strings.Split(strings.TrimSpace(string(slavePorts)), "\n")

		for index := 0; index <= (len(slaveIP) - 1); index++ {
			slaveRoleCommand := ("redis-cli -h " + string(slaveIP[index]) + " -p " + string(slavePort[index]) + " info replication | grep role")
			slaveRole, _ := exec.Command("bash", "-c", slaveRoleCommand).Output()
			slaveRoleSlice = append(slaveRoleSlice, string(slaveRole))

			if strings.TrimSpace(slaveRoleSlice[0]) != "role:slave" {
				fmt.Printf("\n❌  Role in Slave %s:%s is not Slave ", string(slaveIP[index]), string(slavePort[index]))
				replicationFlag = true
			}
		}
	}
	if replicationFlag == false {
		fmt.Println(aurora.Green("\n✅   All Master and Slave roles are correct."))
	}
	fmt.Println()
}

// getting the directory where dump files are located and copying all the persistence files into the new directory
func backupRDBFiles(ver string) {
	var (
		result bool
	)
	fmt.Println("\nCopying RDB files into /data/upgrade/persistence/redis_V" + ver + "/\n")

	// deletes all files and diretories on /var/lib/redis except .rdb files
	delNoRDBFilesCommand := ("sudo find /var/lib/redis/* \\! -name '*.rdb' -delete")
	_, errDel := exec.Command("bash", "-c", delNoRDBFilesCommand).Output()
	if errDel == nil {
		fmt.Println(aurora.Green("\nAll no .rdb files and directories were deleted under /var/lib/redis)"))
	} else {
		fmt.Println(aurora.Red("\nError when trying to delete all no .rdb files under /var/lib/redis"))
		fmt.Println("Press any key if want to continue")
		fmt.Scanln()
	}

	getRDBFilesCopyCommand := ("sudo cp -rf  /var/lib/redis/dump*.rdb /data/upgrade/persistence/redis_V" + ver + "/")
	fmt.Println(getRDBFilesCopyCommand)
	_, err := exec.Command("bash", "-c", getRDBFilesCopyCommand).Output()
	if err != nil {
		fmt.Println(aurora.Red("\n ❌  Couldn't copy persistence files into /data/upgrade/persistence/redis_V" + ver + "/"))
	} else {
		root := "/data/upgrade/persistence/redis_V" + ver + "/"
		fmt.Println(root)
		for {
			result = checkOriginDestinationFiles(root)
			if result == true {
				break
			}
		}
	}
	return
}

//run BGSAVE on each instance
func runBackup(wg *sync.WaitGroup, port string) {
	fmt.Printf("\n*************************  Started BGSAVE on %s  *************************\n", port)
	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:" + port,
		MinIdleConns: 1,
		MaxRetries:   5,
		ReadTimeout:  -1,
		WriteTimeout: -1,
	})

	_, err := client.Save().Result()
	if err != nil {
		fmt.Printf("\nCould not save: %s", err)
	}
	defer wg.Done()
}

// calls runBackup function for each redis port and runs BGSAVE
func runBGSAVEFull() {
	var dataPortsMaster []string
	var wg sync.WaitGroup

	INPUTFILE := "/tmp/redis-upgrade-configuration.txt"
	_, _, _, dataPortsMaster, _, _ = getServersPorts(INPUTFILE)

	for _, port := range dataPortsMaster {
		wg.Add(1)
		go runBackup(&wg, port)
	}
	wg.Wait()
	fmt.Println(aurora.Green("\n✅  BGSAVE completed."))
}

/* In case of the running redis was installed using binaries, the new apt-get installation changed
the redis-server executable file location and supervisor files need to be updated */
func updateSupervisorConfFiles() {
	pattern := []string{"redis*"}
	targetDir := "/etc/supervisor/conf.d/"

	for _, v := range pattern {
		matches, err := filepath.Glob(targetDir + v)
		if err != nil {
			fmt.Println(aurora.Red("\n ❌  None of the files under /etc/supervisor/conf.d starts with redis !"))
		}
		for _, file := range matches {
			input, err := ioutil.ReadFile(file)
			output := bytes.Replace(input, []byte("/usr/local/bin/redis-server"), []byte("/usr/bin/redis-server"), -1)

			if err = ioutil.WriteFile(file, output, 0666); err != nil {
				fmt.Println("\n ❌  Couldn't replace redis-server string on Supervisor conf files: ", err)

			}
		}
	}
}

//preBgsaveCheck checks masters-slaves are synced

func preBgsaveCheckFull() {

	var (
		masterRoleSlice      []string
		slaveRoleSlice       []string
		slaveLinkSlice       []string
		dataPortsMaster      []string
		dataPortsSlaves      []string
		dataPortsSentinel    []string
		slaveServer          []string
		sentinelServer       []string
		clientsNameListSlice []string
		sentinelFlag         = false
		replicationFlag      = false
		keyspaceFlag         = false
		pingResult           = false
		server               string
		instance             string
	)

	fmt.Println("\n*************************  PRE_BGSAVE CHECK  *************************")

	_, sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel = getServersPorts(INPUTFILE)

	// check replication status
	// on master the role should be "master"
	fmt.Println("\n\nChecking replication status:")
	for _, portM := range dataPortsMaster {
		getMasterReplicationRoleCommand := ("redis-cli -p " + portM + " info replication | grep role")
		masterRole, _ := exec.Command("bash", "-c", getMasterReplicationRoleCommand).Output()
		masterRoleSlice = append(masterRoleSlice, string(masterRole))
	}
	for i := 0; i <= (len(masterRoleSlice) - 1); i++ {
		if strings.TrimSpace(masterRoleSlice[i]) != "role:master" {
			fmt.Printf("\n❌  Replication error: role on instance %s is not Master", dataPortsMaster[i])
			replicationFlag = true
		}
	}
	// on slaves the role should be "slave" and also link status "up"

	for n, portS := range dataPortsSlaves {
		getSlaveReplicationRoleCommand := ("redis-cli -h " + string(slaveServer[n]) + " -p " + portS + " info replication | grep role")
		slaveRole, _ := exec.Command("bash", "-c", getSlaveReplicationRoleCommand).Output()
		slaveRoleSlice = append(slaveRoleSlice, string(slaveRole))
	}
	for i := 0; i <= (len(slaveRoleSlice) - 1); i++ {
		if strings.TrimSpace(slaveRoleSlice[i]) != "role:slave" {
			fmt.Printf("\n❌  Replication error: role on Slave server %s instance %s is not Slave", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}
	}
	for r, portS := range dataPortsSlaves {
		getSlaveReplicationLinkCommand := ("redis-cli -h " + string(slaveServer[r]) + " -p " + portS + " info replication | grep ^master_link_status | cut -d \":\" -f 2")
		slaveLink, _ := exec.Command("bash", "-c", getSlaveReplicationLinkCommand).Output()
		slaveLinkSlice = append(slaveLinkSlice, string(slaveLink))
	}
	for i := 0; i <= (len(slaveLinkSlice) - 1); i++ {
		if strings.TrimSpace(slaveLinkSlice[i]) != "up" {
			fmt.Printf("\n❌  Replication error on Slave server %s: link on instance %s is not \"up\"", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}

	}
	if replicationFlag == false {
		fmt.Println(aurora.Green("\n✅  Master and Slave roles are correct."))
		fmt.Println(aurora.Green("\n✅  Master-Slave link is \"up\"."))
	}
	fmt.Println()

	//check PING
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
			fmt.Println(aurora.Red("❌  Check again later. Leaving the script."))
			return
		}
	}
	//checks min_slaves_good_slaves
	for _, portM := range dataPortsMaster {
		getMinGoodSlavesCommand := ("redis-cli -p " + portM + " info replication | grep connected_slaves | cut -d \":\" -f 2")
		getMinGoodSlaves, _ := exec.Command("bash", "-c", getMinGoodSlavesCommand).Output()

		if strings.TrimSpace(string(getMinGoodSlaves)) != "2" {
			fmt.Printf("\n❌ Number of Slaves linked to Master Port %s is %s \n", portM, getMinGoodSlaves)
		}
	}

	//checking Sentinel status
	sentinelFlag = checkSentinelStatus()
	if sentinelFlag == false {
		fmt.Println(aurora.Green("\n✅  Sentinel status are all OK."))
	}

	// check clients connected
	fmt.Println("\n\nChecking if there are clients still connected:")
	for _, portM := range dataPortsMaster {
		//		clientsNameListSlice = nil

		getClientsCommand := ("redis-cli -p " + portM + " client list|awk {'print $2'} |awk -F= {'print $2'}|awk -F: {'print $1'} |sort |uniq")
		clientsList, _ := exec.Command("bash", "-c", getClientsCommand).Output()
		clientListSlice := strings.Split(strings.TrimSpace(string(clientsList)), "\n")

		for index := 0; index <= (len(clientListSlice) - 1); index++ {
			getClientNameCommand := ("nslookup " + string(clientListSlice[index]) + " | grep name |awk {'print $4'}")
			clientsNameList, _ := exec.Command("bash", "-c", getClientNameCommand).Output()
			clientsNameListSlice = append(clientsNameListSlice, string(clientsNameList))
		}

		fmt.Printf("\n✅  List of clients still connected to instance %s is:\n", portM)
		fmt.Println(strings.Trim(fmt.Sprint(clientsNameListSlice), "[]"))
	}

	// get the Master number of keys and compare. If they don't match, raise an error
	fmt.Println("\nChecking Master-Slave number of keys matching:")

	for _, portM := range dataPortsMaster {
		masterKeysCommand := ("redis-cli -p " + portM + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
		masterKeys, _ := exec.Command("bash", "-c", masterKeysCommand).Output()

		masterSlaveIPCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 2 | cut -d \",\" -f 1")
		slaveIPs, _ := exec.Command("bash", "-c", masterSlaveIPCommand).Output()

		masterSlavePortCommand := ("redis-cli -p " + portM + " info  replication | grep ^slave | cut -d \"=\" -f 3 | cut -d \",\" -f 1")
		slavePorts, _ := exec.Command("bash", "-c", masterSlavePortCommand).Output()

		slaveIP := strings.Split(strings.TrimSpace(string(slaveIPs)), "\n")
		slavePort := strings.Split(strings.TrimSpace(string(slavePorts)), "\n")

		for index := 0; index <= (len(slaveIP) - 1); index++ {
			slaveKeysCommand := ("redis-cli -h " + string(slaveIP[index]) + " -p " + string(slavePort[index]) + " info Keyspace | grep db | cut -d '=' -f 2 | cut -d ',' -f 1")
			slaveKeys, _ := exec.Command("bash", "-c", slaveKeysCommand).Output()
			if string(masterKeys) != string(slaveKeys) {
				keyspaceFlag = true
				fmt.Printf(" ❌  Keys don't match on Master port %s - Slave %s:%s \n", portM, string(slaveIP[index]), string(slavePort[index]))
			}
		}
	}
	if keyspaceFlag == false {
		fmt.Println(aurora.Green("\n✅   Number of KEYS on Master - Slaves are all OK"))
	}

}

// get Master, Slave, Sentinel slices from input file

func getServersPorts(INPUTFILE string) ([]string, []string, []string, []string, []string, []string) {

	var (
		dataPortsMasterSlice    []string
		sentinelPortSlice       []string
		sentinelServerNameSlice []string
		slaveServerIPSlice      []string
		slaveServerPortSlice    []string
		masterServerNameSlice   []string
	)

	//reads the file line by line and getting the servers and ports
	file, err := os.Open(INPUTFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineArray := strings.Fields(scanner.Text())
		redisPort := lineArray[0]
		redisHost := lineArray[1]
		redisRole := lineArray[2]

		//appending the server name, portSlice slice with all the ports detailed on INPUTFILE

		if redisRole == "master" {
			dataPortsMasterSlice = append(dataPortsMasterSlice, redisPort)
			masterServerNameSlice = append(masterServerNameSlice, redisHost)
		} else if redisRole == "sentinel" {
			sentinelPortSlice = append(sentinelPortSlice, redisPort)
			sentinelServerNameSlice = append(sentinelServerNameSlice, redisHost)
		} else if redisRole == "slave" {
			slaveServerPortSlice = append(slaveServerPortSlice, redisPort)
			slaveServerIPSlice = append(slaveServerIPSlice, redisHost)
		}
	}

	return masterServerNameSlice, sentinelServerNameSlice, slaveServerIPSlice, dataPortsMasterSlice, slaveServerPortSlice, sentinelPortSlice
}

//postBgsaveCheck checks bgsave completed correctly

func postBgsaveCheckFull() {

	var (
		dataPortsMaster   []string
		dataPortsSlaves   []string
		dataPortsSentinel []string
		slaveServer       []string
		sentinelServer    []string
		pingResult        = false
		server            string
		instance          string
		masterRoleSlice   []string
		slaveRoleSlice    []string
		replicationFlag   = false
		slaveLinkSlice    []string
	)

	_, sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel = getServersPorts(INPUTFILE)

	fmt.Println("\n*************************  POST_BGSAVE CHECK  *************************")

	// check replication status
	// on master the role should be "master"
	fmt.Println("\n\nChecking replication status:")
	for _, portM := range dataPortsMaster {
		getMasterReplicationRoleCommand := ("redis-cli -p " + portM + " info replication | grep role")
		masterRole, _ := exec.Command("bash", "-c", getMasterReplicationRoleCommand).Output()
		masterRoleSlice = append(masterRoleSlice, string(masterRole))
	}
	for i := 0; i <= (len(masterRoleSlice) - 1); i++ {
		if strings.TrimSpace(masterRoleSlice[i]) != "role:master" {
			fmt.Printf("\n❌  Replication error: role on instance %s is not Master", dataPortsMaster[i])
			replicationFlag = true
		}
	}
	// on slaves the role should be "slave" and also link status "up"

	for n, portS := range dataPortsSlaves {
		getSlaveReplicationRoleCommand := ("redis-cli -h " + string(slaveServer[n]) + " -p " + portS + " info replication | grep role")
		slaveRole, _ := exec.Command("bash", "-c", getSlaveReplicationRoleCommand).Output()
		slaveRoleSlice = append(slaveRoleSlice, string(slaveRole))
	}
	for i := 0; i <= (len(slaveRoleSlice) - 1); i++ {
		if strings.TrimSpace(slaveRoleSlice[i]) != "role:slave" {
			fmt.Printf("\n❌  Replication error: role on Slave server %s instance %s is not Slave", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}
	}
	for r, portS := range dataPortsSlaves {
		getSlaveReplicationLinkCommand := ("redis-cli -h " + string(slaveServer[r]) + " -p " + portS + " info replication | grep ^master_link_status | cut -d \":\" -f 2")
		slaveLink, _ := exec.Command("bash", "-c", getSlaveReplicationLinkCommand).Output()
		slaveLinkSlice = append(slaveLinkSlice, string(slaveLink))
	}
	for i := 0; i <= (len(slaveLinkSlice) - 1); i++ {
		if strings.TrimSpace(slaveLinkSlice[i]) != "up" {
			fmt.Printf("\n❌  Replication error on Slave server %s: link on instance %s is not \"up\"", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}

	}
	if replicationFlag == false {
		fmt.Println(aurora.Green("\n✅  Master and Slave roles are correct."))
		fmt.Println(aurora.Green("\n✅  Master-Slave link is \"up\"."))
	}
	fmt.Println()

	//check PING
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
			fmt.Println(aurora.Red("❌  Check again later. Leaving the script."))
			return
		}
	}
	//checks min_slaves_good_slaves
	for _, portM := range dataPortsMaster {
		getMinGoodSlavesCommand := ("redis-cli -p " + portM + " info replication | grep connected_slaves | cut -d \":\" -f 2")
		getMinGoodSlaves, _ := exec.Command("bash", "-c", getMinGoodSlavesCommand).Output()

		if strings.TrimSpace(string(getMinGoodSlaves)) != "2" {
			fmt.Printf("\n❌ Number of Slaves linked to Master Port %s is %s \n", portM, getMinGoodSlaves)
		}
	}

	// check last bgsave
	fmt.Println("\n\nChecking last BGsave time:")
	for _, portM := range dataPortsMaster {
		getLastsaveCommand := ("redis-cli -p " + portM + " lastsave")
		getLastsave, errLastsaveM := exec.Command("bash", "-c", getLastsaveCommand).Output()
		if errLastsaveM != nil {
			fmt.Printf("\n❌  Unable to get LAST BGSAVE on localhost instance %s\n", portM)
		}
		timeNow := time.Now().Unix()
		fmt.Printf("Instance %s: \tLast Save: %s\t\t Time now: %d\n", portM, getLastsave, timeNow)
		fmt.Println("------------------------------------------")
	}

	//	checks if bgsaving process is still running. It shouldn't
	fmt.Println("\n\nChecking if BGsave process still running:")
	output, _ := exec.Command("ps -aef | grep -i bgsave | grep -v color").Output()
	if string(output) != "" {
		fmt.Printf("\n❌  BGSAVE process is still running! \n  %s", string(output))
	} else {
		fmt.Println(aurora.Green("\n✅  No BGsave process still running."))
	}
	fmt.Println()
}

func checkPing(sentinelServer, slaveServer, portsMaster, portsSlaves, portsSentinel []string) (bool, string, string) {
	var failPing = false

	for _, portM := range portsMaster {
		getPingCommand := ("redis-cli -p " + portM + " ping")
		_, errPingM := exec.Command("bash", "-c", getPingCommand).Output()
		if errPingM != nil {
			failPing = true
			return failPing, "localhost", portM
		}
	}

	for index, portS := range portsSlaves {
		getPingCommand := ("redis-cli -h " + strings.TrimSpace(string(slaveServer[index])) + " -p " + portS + " ping")
		_, errPingS := exec.Command("bash", "-c", getPingCommand).Output()
		if errPingS != nil {
			failPing = true
			return failPing, slaveServer[index], portS
		}
	}
	for index, portSent := range portsSentinel {
		getPingCommand := ("redis-cli -h " + sentinelServer[index] + " -p " + portSent + " ping")
		_, errPingS := exec.Command("bash", "-c", getPingCommand).Output()
		if errPingS != nil {
			failPing = true
			return failPing, sentinelServer[index], portSent
		}
	}

	return false, " ", " "
}

// On Slaves: pre-isntallation check

func preInstallCheckSlave() {
	var (
		slaveRoleSlice    []string
		slaveLinkSlice    []string
		dataPortsMaster   []string
		dataPortsSlaves   []string
		dataPortsSentinel []string
		slaveServer       []string
		sentinelServer    []string
		masterServer      []string
		sentinelFlag      = false
		replicationFlag   = false
		pingResult        = false
		server            string
		instance          string
	)

	fmt.Println("\n*************************  PRE INSTALLATION CHECK - SLAVES  *************************")

	masterServer, sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel = getServersPorts(INPUTFILE)

	// on slaves the role should be "slave" and also link status "up"

	for n, portS := range dataPortsSlaves {
		getSlaveReplicationRoleCommand := ("redis-cli -h " + string(slaveServer[n]) + " -p " + portS + " info replication | grep role")
		slaveRole, _ := exec.Command("bash", "-c", getSlaveReplicationRoleCommand).Output()
		slaveRoleSlice = append(slaveRoleSlice, string(slaveRole))
	}
	for i := 0; i <= (len(slaveRoleSlice) - 1); i++ {
		if strings.TrimSpace(slaveRoleSlice[i]) != "role:slave" {
			fmt.Printf("\n❌  Replication error: role on Slave server %s instance %s is not Slave", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}
	}
	for r, portS := range dataPortsSlaves {
		getSlaveReplicationLinkCommand := ("redis-cli -h " + string(slaveServer[r]) + " -p " + portS + " info replication | grep ^master_link_status | cut -d \":\" -f 2")
		slaveLink, _ := exec.Command("bash", "-c", getSlaveReplicationLinkCommand).Output()
		slaveLinkSlice = append(slaveLinkSlice, string(slaveLink))
	}
	for i := 0; i <= (len(slaveLinkSlice) - 1); i++ {
		if strings.TrimSpace(slaveLinkSlice[i]) != "up" {
			fmt.Printf("\n❌  Replication error on Slave server %s: link on instance %s is not \"up\"", string(slaveServer[i]), dataPortsSlaves[i])
			replicationFlag = true
		}

	}
	if replicationFlag == false {
		fmt.Println(aurora.Green("\n✅  Master and Slave roles are correct."))
		fmt.Println(aurora.Green("\n✅  Master-Slave link is \"up\"."))
	}
	fmt.Println()

	//check PING
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
			fmt.Println(aurora.Red("❌  Check again later. Leaving the script."))
			return
		}
	}
	//checks min_slaves_good_slaves
	for t, portM := range dataPortsMaster {
		getMinGoodSlavesCommand := ("redis-cli -h " + masterServer[t] + " -p " + portM + " info replication | grep connected_slaves | cut -d \":\" -f 2")
		getMinGoodSlaves, _ := exec.Command("bash", "-c", getMinGoodSlavesCommand).Output()

		if strings.TrimSpace(string(getMinGoodSlaves)) != "2" {
			fmt.Printf("\n❌ Number of Slaves linked on Master Server %s port %s is %s \n", masterServer[t], portM, getMinGoodSlaves)

		}
	}

	//checking Sentinel status
	sentinelFlag = checkSentinelStatus()
	if sentinelFlag == false {
		fmt.Println(aurora.Green("\n✅  Sentinel status are all OK."))
	}

}

func checkSentinel() {
	var (
		dataPortsMaster   []string
		dataPortsSlaves   []string
		dataPortsSentinel []string
		slaveServer       []string
		sentinelServer    []string
		sentinelFlag      = false
		pingResult        = false
		server            string
		instance          string
	)

	fmt.Println("\n*************************  SENTINEL CHECK  *************************")

	_, sentinelServer, _, _, _, dataPortsSentinel = getServersPorts(INPUTFILE)

	// on slaves the role should be "slave" and also link status "up"

	//check PING
	pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
	if pingResult == true {
		fmt.Printf("\n❌  Ping failed on server %s instance %s\n", server, instance)
		fmt.Printf("Checking again in 5 seconds...")
		time.Sleep(5 * time.Second)
		pingResult, server, instance = checkPing(sentinelServer, slaveServer, dataPortsMaster, dataPortsSlaves, dataPortsSentinel)
		if pingResult == true {
			fmt.Printf("\n❌  Ping failed again on server %s instance %s\n", server, instance)
			fmt.Println(aurora.Red("❌  Check again later. Leaving the script."))
			return
		}
	}
	fmt.Println(aurora.Green("\n✅  Ping Sentinels are all OK."))

	//checking Sentinel status
	sentinelFlag = checkSentinelStatus()
	if sentinelFlag == false {
		fmt.Println(aurora.Green("\n✅  Sentinel status are all OK."))
	}

}

// standing on Master server, checks if every Redis instance is master
func checkMasterRole() {
	var (
		dataPortsMaster []string
		masterRoleSlice []string
	)

	_, _, _, dataPortsMaster, _, _ = getServersPorts(INPUTFILE)

	fmt.Println("\n\nChecking Master instance status:")
	for _, portM := range dataPortsMaster {
		getMasterReplicationRoleCommand := ("redis-cli -p " + portM + " info replication | grep role")
		masterRole, _ := exec.Command("bash", "-c", getMasterReplicationRoleCommand).Output()
		masterRoleSlice = append(masterRoleSlice, string(masterRole))
	}
	for i := 0; i <= (len(masterRoleSlice) - 1); i++ {
		if strings.TrimSpace(masterRoleSlice[i]) != "role:master" {
			fmt.Printf("\n❌  Replication error: role on instance %s is not Master", dataPortsMaster[i])
		}
	}
}

// menu options
func menu() {
	var input int
	var verShortNumber string

	n, err := fmt.Scanln(&input)
	if n < 1 || err != nil {
		fmt.Println(aurora.Red("\n❌  Invalid option."))
		return
	}

	switch input {
	case 1:
		preBgsaveCheckFull()
		fmt.Printf("\nWhat's next? ")
	case 2:
		runBGSAVEFull()
		fmt.Printf("\nWhat's next? ")
	case 3:
		postBgsaveCheckFull()
		fmt.Printf("\nWhat's next? ")
	case 4:
		verShortNumber, _ := versionToCopy()
		dataDirCreation(verShortNumber)
		backupRDBFiles(verShortNumber)
		backupConfFiles()
		backupSupervisorFiles()
		fmt.Printf("\nWhat's next? ")
	case 5:
		verShortNumber, versionToInstall := versionToInstallFunc()
		uninstallAndInstallRedis(versionToInstall)
		currentVersion()
		stopRedis6379(6379)
		disableSystemctlRedis()
		updateSupervisorConfFiles()
		reloadSupervisor()
		addProtectedModeNoMaster()
		copyingBackRDBFiles(verShortNumber)
		reloadSupervisor()
		checkDirOnConfFiles()
		checkMasterRole()
		fmt.Printf("\nWhat's next? ")
	case 6:
		finalCheckMaster()
		fmt.Printf("\nWhat's next? ")
	case 10:
		preInstallCheckSlave()
		fmt.Printf("\nWhat's next? ")
	case 11:
		verShortNumber, versionToInstall := versionToInstallFunc()
		dataDirCreation(verShortNumber)
		backupConfFiles()
		backupSupervisorFiles()
		uninstallAndInstallRedis(versionToInstall)
		currentVersion()
		fmt.Printf("\nWhat's next? ")
	case 12:
		addProtectedModeNoSentinelsOnSlaves(verShortNumber)
		reloadSupervisor()
		fmt.Printf("\nWhat's next? ")
	case 13:
		finalCheckSlaves()
		fmt.Printf("\nWhat's next? ")
	case 21:
		checkSentinel()
		fmt.Printf("\nWhat's next? ")
	case 22:
		verShortNumber, versionToInstall := versionToInstallFunc()
		dataDirCreation(verShortNumber)
		backupConfFiles()
		backupSupervisorFiles()
		uninstallAndInstallRedis(versionToInstall)
		currentVersion()
		fmt.Printf("\nWhat's next? ")
	case 23:
		addProtectedModeNoSentinelsOnSentinelServers()
		fmt.Printf("\nWhat's next? ")
	case 24:
		checkSentinel()
		fmt.Printf("\nWhat's next? ")
	case 212:
		var input2 int
		var verShortNumber string

		fmt.Println("\nPress 101 for uninstall And Install Redis().")
		fmt.Println("\nPress 102 for copying rdb files info /data/upgrade/peristence/...")
		fmt.Println("\nPress 103 for stop And Disable Redis 6379().")
		fmt.Println("\nPress 104 for update Supervisor Conf Files().")
		fmt.Println("\nPress 105 for reload Supervisor().")
		fmt.Println("\nPress 106 for add Protected Mode No Master().")
		fmt.Println("\nPress 107 for add Protected Mode No Slaves().")
		fmt.Println("\nPress 108 for add Protected Mode No Sentinels On Sentinel Servers().")
		fmt.Println("\nPress 109 for copying Back RDB Files().")
		fmt.Println("\nPress 110 for check Dir Value on Conf Files().")
		fmt.Println("\nPress 111 for check Master Role().")

		n, err2 := fmt.Scanln(&input2)
		if n < 1 || err2 != nil {
			fmt.Println(aurora.Red("\n❌  Invalid option."))
			return
		}

		switch input2 {
		case 101:
			_, versionToInstall := versionToInstallFunc()
			uninstallAndInstallRedis(versionToInstall)
			fmt.Printf("\nWhat's next? ")
		case 102:
			verShortNumber, _ := versionToCopy()
			backupRDBFiles(verShortNumber)
			fmt.Printf("\nWhat's next? ")
		case 103:
			stopRedis6379(6379)
			disableSystemctlRedis()
			fmt.Printf("\nWhat's next? ")
		case 104:
			updateSupervisorConfFiles()
			fmt.Printf("\nWhat's next? ")
		case 105:
			reloadSupervisor()
			fmt.Printf("\nWhat's next? ")
		case 106:
			addProtectedModeNoMaster()
			fmt.Printf("\nWhat's next? ")
		case 107:
			addProtectedModeNoSentinelsOnSlaves(verShortNumber)
			fmt.Printf("\nWhat's next? ")
		case 108:
			addProtectedModeNoSentinelsOnSentinelServers()
			fmt.Printf("\nWhat's next? ")
		case 109:
			copyingBackRDBFiles(verShortNumber)
			fmt.Printf("\nWhat's next? ")
		case 110:
			checkDirOnConfFiles()
			fmt.Printf("\nWhat's next? ")
		case 111:
			checkMasterRole()
			fmt.Printf("\nWhat's next? ")
		}

		fmt.Printf("\nWhat's next? ")
	}
}

// main function

func main() {

	fileStat, err := os.Stat("/tmp/redis-upgrade-configuration.txt")

	if err != nil {
		fmt.Println(aurora.Red("Could not get /tmp/redis-upgrade-configuration.txt file"))
		log.Fatal(err)
	}
	if fileStat.Size() == 0 {
		fmt.Println(aurora.Red("/tmp/redis-upgrade-configuration.txt file is emtpty."))
		log.Fatal(err)
	}
	currentVersion()
	fmt.Println()
	fmt.Println(aurora.Bold("\n  👍  ON MASTER SERVER  👍  "))
	fmt.Println("\nPress 1 for Pre BGSAVE Check.")
	fmt.Println("\nPress 2 for run BGSAVE.")
	fmt.Println("\nPress 3 for Post BGSAVE Check.")
	fmt.Println("\nPress 4 for Backup RDB Files.")
	fmt.Println("\nPress 5 for Uninstall/Install Redis. Copying back RDB files.  Adding \"Protected Mode no\". Checking Master role.")
	fmt.Println("\nPress 6 for general final check.")
	fmt.Println(aurora.Bold("\n  👍  ON SLAVE SERVERS  👍  "))
	fmt.Println("\nPress 10 for pre-install check.")
	fmt.Println("\nPress 11 for Uninstall/Install Redis.")
	fmt.Println("\nPress 12 for adding \"Protected Mode no\" on Slaves.")
	fmt.Println("\nPress 13 for Post upgrade check.")
	fmt.Println(aurora.Bold("\n  👍  ON SENTINEL ONLY SERVERS  👍  "))
	fmt.Println("\nPress 21 for pre-install check.")
	fmt.Println("\nPress 22 for Uninstall/Install Redis.")
	fmt.Println("\nPress 23 for adding \"Protected Mode no\" on Sentinels.")
	fmt.Println("\nPress 24 for Post upgrade check.")
	fmt.Println()
	fmt.Println()
	fmt.Println("\nPress 212 - DEBUG mode")
	fmt.Println()
	for {
		menu()
	}
}
