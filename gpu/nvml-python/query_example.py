#pynvml github, implemenation ref https://github.com/gpuopenanalytics/pynvml/blob/master/pynvml/nvml.py
from pynvml import *

def get_utilization_per_process():
    nvmlInit()
    deviceCount = nvmlDeviceGetCount()
    for i in range(deviceCount):
        handle = nvmlDeviceGetHandleByIndex(i)
        use = nvmlDeviceGetUtilizationRates(handle)
        print("Device {}, GPU util rate {}, Memory util Rate {}".format(i, use.gpu, use.memory))
        if use.gpu > 0:
            print("Device {} is used".format(i))
            # get memory used per process
            processes = nvmlDeviceGetComputeRunningProcesses(handle)
            if len(processes) > 0:
                for j in processes:
                    print("process id {}, used memory {}".format(j.pid,j.usedGpuMemory))

            # get gpu util used per process
            lastseentimestamp = 0 # set 0 to return all samples buffered
            samples = nvmlDeviceGetProcessUtilization(handle, lastseentimestamp)
            for sample in samples:
                #available attribute 'decUtil', 'encUtil', 'memUtil', 'pid', 'smUtil', 'timeStamp'
                print("process id {}, timestamp {}, sm util {}, mem util {}, decUtil {}, encUtil:{}".format(sample.pid, sample.timeStamp,sample.smUtil,sample.memUtil,sample.decUtil,sample.encUtil))

        else:
            print("Device {} is free".format(i))


if __name__ == "__main__":
   get_utilization_per_process()
