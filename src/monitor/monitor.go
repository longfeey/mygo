package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	/* CPU */
	mtCPULl           = "/proc/cpufreq/MT_CPU_DVFS_LL"
	mtCPUL            = "/proc/cpufreq/MT_CPU_DVFS_L"
	cpuFreq           = "cpufreq_freq"
	cpuOnline         = "/sys/devices/system/cpu/online"
	cpu0CurFreq       = "/sys/devices/system/cpu/cpufreq/policy0/scaling_cur_freq"
	cpu4CurFreq       = "/sys/devices/system/cpu/cpufreq/policy4/scaling_cur_freq"
	cpu0AvailableFreq = "/sys/devices/system/cpu/cpufreq/policy0/scaling_available_frequencies"
	cpu4AvailableFreq = "/sys/devices/system/cpu/cpufreq/policy4/scaling_available_frequencies"
	cpu0TimeInState   = "/sys/devices/system/cpu/cpufreq/policy0/stats/time_in_state"
	cpu4TimeInState   = "/sys/devices/system/cpu/cpufreq/policy4/stats/time_in_state"
	cpuInfo           = "/proc/cpuinfo"

	/* GPU */
	gpuPowerDump = "/proc/gpufreq/gpufreq_power_dump"
)

var cpuNum uint64

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	// Stdout pipe for reading the generated output.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the script.
	if err := cmd.Run(); err != nil {
		c := name
		if len(args) > 0 {
			c += " " + strings.Join(args, " ")
		}
		return "", fmt.Errorf("failed to run command %q:\n  %v\n  %s", c, err, stderr.String())
	}

	return stdout.String(), nil
}

func cpuinfoGet() {
	/*Get cpu0~cpu3 current frequency */
	cpu0freq, _ := runCommand("adb", "shell", "cat", cpu0CurFreq)
	cpu0freq = strings.TrimSpace(cpu0freq)

	/*Get cpu online */
	cpuOnlineNum, _ := runCommand("adb", "shell", "cat", cpuOnline)
	cpuOnlineNum = strings.TrimSpace(cpuOnlineNum)

	if cpuNum == 8 {
		/*Get cpu4~cpu7 current frequency */
		cpu4freq, _ := runCommand("adb", "shell", "cat", cpu4CurFreq)
		cpu4freq = strings.TrimSpace(cpu4freq)
		fmt.Printf("%-8s\t%-8s\t%s\n", cpu0freq, cpu4freq, cpuOnlineNum)
	} else if cpuNum == 4 {
		fmt.Printf("%-8s\t%s\n", cpu0freq, cpuOnlineNum)
	}
}

func cpuAvailableFreqGet() {
	/*Get cpu0~cpu3 available frequency */
	cpu0AvailableFreq, _ := runCommand("adb", "shell", "cat", cpu0AvailableFreq)
	//cpu0AvailableFreq = strings.TrimSpace(cpu0AvailableFreq)
	fmt.Printf("cpu0~3 available frequencies:\n%s\n", cpu0AvailableFreq)

	if cpuNum == 8 {
		/*Get cpu4~cpu7 available frequency */
		cpu4AvailableFreq, _ := runCommand("adb", "shell", "cat", cpu4AvailableFreq)
		//cpu4available_freq = strings.TrimSpace(cpu4available_freq)
		fmt.Printf("cpu4~7 available_frequencies:\n%s\n", cpu4AvailableFreq)
	}
	fmt.Printf("\n")
}

func gpuPowerDumpGet() {
	/*Get Gpu available frequency */
	gpuPowerDumpInfo, _ := runCommand("adb", "shell", "cat", gpuPowerDump)
	//cpu0available_freq = strings.TrimSpace(cpu0available_freq)
	fmt.Printf("GPU power dump:\n%s\n", gpuPowerDumpInfo)
}

func cpuNumGet() uint64 {
	num, err := runCommand("adb", "shell", "cat", cpuInfo, "|", "grep", "processor", "|", "wc", "-l")
	if err != nil {
		fmt.Printf("get cpu num failed,error:%s", err)
		return 0
	}
	fmt.Printf("get cpu num %s", num)
	//fmt.Printf("get cpu num %d",strconv.Atoi(num))
	num = strings.TrimSpace(num)

	ret, err := strconv.ParseUint(num, 10, 0)
	if err != nil {
		panic(err)
	}

	return ret
}

func main() {
	cpuNum = cpuNumGet()

	gpuPowerDumpGet()

	cpuAvailableFreqGet()
	if cpuNum == 8 {
		fmt.Printf("%-8s\t%-8s\t%s\n", "CPU0~CPU3", "CPU4~CPU7", "CPU online")
	} else {
		fmt.Printf("%-8s\t%s\n", "CPU0~CPU3", "CPU online")
	}

	for {
		cpuinfoGet()
	}
}
