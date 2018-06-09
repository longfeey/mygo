package main

import (
    "fmt"
    "bytes"
    "strings"
    "os/exec"
    "strconv"
)

const (
    /* CPU */
    mt_cpu_ll = "/proc/cpufreq/MT_CPU_DVFS_LL"
    mt_cpu_l = "/proc/cpufreq/MT_CPU_DVFS_L"
    cpufreq = "cpufreq_freq"
    cpu_online = "/sys/devices/system/cpu/online"
    cpu0_cur_freq = "/sys/devices/system/cpu/cpufreq/policy0/scaling_cur_freq"
    cpu4_cur_freq = "/sys/devices/system/cpu/cpufreq/policy4/scaling_cur_freq"
    cpu0_available_freq = "/sys/devices/system/cpu/cpufreq/policy0/scaling_available_frequencies"
    cpu4_available_freq = "/sys/devices/system/cpu/cpufreq/policy4/scaling_available_frequencies"
    cpu0_time_in_state = "/sys/devices/system/cpu/cpufreq/policy0/stats/time_in_state"
    cpu4_time_in_state = "/sys/devices/system/cpu/cpufreq/policy4/stats/time_in_state"
    cpu_info = "/proc/cpuinfo"

    /* GPU */
    gpu_power_dump = "/proc/gpufreq/gpufreq_power_dump"
)

var cpu_num uint64 

func RunCommand(name string, args ...string) (string, error) {
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

func cpuinfo_get() {
    /*Get cpu0~cpu3 current frequency */
    cpu0freq ,_ := RunCommand("adb", "shell", "cat", cpu0_cur_freq)
    cpu0freq = strings.TrimSpace(cpu0freq)

    /*Get cpu online */
    cpu_online_num,_ := RunCommand("adb", "shell", "cat", cpu_online)
    cpu_online_num = strings.TrimSpace(cpu_online_num)

    if cpu_num == 8 {
        /*Get cpu4~cpu7 current frequency */
        cpu4freq ,_ := RunCommand("adb", "shell", "cat", cpu4_cur_freq)
        cpu4freq = strings.TrimSpace(cpu4freq)
        fmt.Printf("%-8s\t%-8s\t%s\n",cpu0freq,cpu4freq,cpu_online_num)
    } else {
        fmt.Printf("%-8s\t%s\n",cpu0freq,cpu_online_num)
    }
}

func cpu_available_freq_get() {
    /*Get cpu0~cpu3 available frequency */
    cpu0available_freq ,_ := RunCommand("adb", "shell", "cat", cpu0_available_freq)
    //cpu0available_freq = strings.TrimSpace(cpu0available_freq)
    fmt.Printf("cpu0~3 available_frequencies:\n%s\n",cpu0available_freq)

    if cpu_num == 8 {
        /*Get cpu4~cpu7 available frequency */
        cpu4available_freq ,_ := RunCommand("adb", "shell", "cat", cpu4_available_freq)
        //cpu4available_freq = strings.TrimSpace(cpu4available_freq)
        fmt.Printf("cpu4~7 available_frequencies:\n%s\n",cpu4available_freq)
    }
    fmt.Printf("\n")
}

func gpu_power_dump_get() {
    /*Get Gpu available frequency */
    gpu_power_dump_info,_ := RunCommand("adb", "shell", "cat", gpu_power_dump)
    //cpu0available_freq = strings.TrimSpace(cpu0available_freq)
    fmt.Printf("GPU power dump:\n%s\n",gpu_power_dump_info)
}

func cpu_num_get() uint64 {
    num, err := RunCommand("adb", "shell", "cat", cpu_info, "|", "grep", "processor", "|", "wc", "-l")
    if err != nil {
        fmt.Printf("get cpu num failed,error:%s",err)
        return 0
    }
    fmt.Printf("get cpu num %s",num)
        //fmt.Printf("get cpu num %d",strconv.Atoi(num))
    num = strings.TrimSpace(num)

    ret, err := strconv.ParseUint(num, 10, 0)
    if err != nil {
        panic(err)
    }

    return ret
}

func main() {
    cpu_num = cpu_num_get()

    gpu_power_dump_get()

    cpu_available_freq_get()
    if cpu_num == 8 {
        fmt.Printf("%-8s\t%-8s\t%s\n","CPU0~CPU3","CPU4~CPU7","CPU online")
    } else {
        fmt.Printf("%-8s\t%s\n","CPU0~CPU3","CPU online")
    }

    for {
        cpuinfo_get()
    }
}
