Install the lasted pynvml package, or install from their [github](https://github.com/gpuopenanalytics/pynvml)

We basically want to query device-level utilization (memory, io, sm) and process level utilization. Important functions are used in the [query-example.py](./nvml-python/query_example.py) file

Sample outputs
```
Device 0, GPU util rate 0, Memory util Rate 0
Device 0 is free
Device 1, GPU util rate 100, Memory util Rate 0
Device 1 is used
process id 4145, used memory 936378368
process id 4145, timestamp 1653167163586231, sm util 99, mem util 0, decUtil 0, encUtil:0
Device 2, GPU util rate 100, Memory util Rate 0
Device 2 is used
process id 4146, used memory 936378368
process id 4146, timestamp 1653167151828549, sm util 99, mem util 0, decUtil 0, encUtil:0
Device 3, GPU util rate 100, Memory util Rate 0
Device 3 is used
process id 4147, used memory 936378368
process id 4147, timestamp 1653167155435855, sm util 99, mem util 0, decUtil 0, encUtil:0
Device 4, GPU util rate 0, Memory util Rate 0
Device 4 is free
Device 5, GPU util rate 0, Memory util Rate 0
Device 5 is free
Device 6, GPU util rate 0, Memory util Rate 0
Device 6 is free
Device 7, GPU util rate 0, Memory util Rate 0
Device 7 is free

```
