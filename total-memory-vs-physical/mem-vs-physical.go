package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

//getDifference gets the difference between the available and the totan redis possible usage
func getDifference(totalMaxmemLimit int, available int) string {
	//	fmt.Printf("Max Mem Limit: %d\n", totalMaxmemLimit)
	//	fmt.Printf("Available: %d\n", available)
	if totalMaxmemLimit > available {
		return "WRONG"
	}
	return "OK"
}

//getMemAvailable gets the total physical memory available
func getMemAvailable() (serverMemAvailableInt int) {
	serverMemAvailableCommand := ("free -b | grep Mem: | awk '{print $2}' ")
	serverMemAvailable, _ := exec.Command("bash", "-c", serverMemAvailableCommand).Output()
	serverMemAvailableInt, _ = strconv.Atoi(strings.TrimSpace(string(serverMemAvailable[:])))
	return serverMemAvailableInt
}

// getAggregation calculates the max mem limit total on the server
func getAggregation(redisMemUsageStringSlice []string) (totalMemUsage int) {
	var (
		ints                  int
		redisMemUsageIntSlice []int
	)

	for i, s := range redisMemUsageStringSlice {
		ints, _ = strconv.Atoi(s)
		redisMemUsageIntSlice = append(redisMemUsageIntSlice, ints)
		totalMemUsage += redisMemUsageIntSlice[i]
	}
	return totalMemUsage

}

// getUsage returns each instance max memroy limit in a []string
func getMemLimits() (r []string) {

	var (
		redisMemUsageStringSlice []string
	)

	redisPortsCommand := ("ps -ef | grep redis-server | grep -v sentinel | grep -v grep | awk 'BEGIN { FS = \"redis-server\" } ; { print $2 }' | grep : | sed 's/*://g'| cut -f2 -d \":\" | sort")
	redisPorts, _ := exec.Command("bash", "-c", redisPortsCommand).Output()
	redisPortsSlice := strings.Split(strings.TrimSpace(string(redisPorts)), "\n")

	for _, port := range redisPortsSlice {
		redisMemUsageCommand := ("redis-cli -p " + port + " config get maxmemory | grep -v maxmemory")
		redisMemUsage, _ := exec.Command("bash", "-c", redisMemUsageCommand).Output()
		redisUsageString := strings.TrimSpace(string(redisMemUsage))
		redisMemUsageStringSlice = append(redisMemUsageStringSlice, redisUsageString)
	}
	return redisMemUsageStringSlice

}

func main() {

	redisMemLimitStringSlice := getMemLimits()
	totalMemLimit := getAggregation(redisMemLimitStringSlice)
	memAvailable := getMemAvailable()
	result := getDifference(totalMemLimit, memAvailable)
	fmt.Println(result)
}
