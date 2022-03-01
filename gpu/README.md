# GPU node preparation

The latest cuda driver 510 only support those [cards](https://docs.nvidia.com/datacenter/tesla/tesla-release-notes-510-47-03/index.html), if you have older card, e.g., K80, install 470.

## 1. Install cuda driver
remove existing one and install the latest
```
sudo apt clean
sudo apt update
sudo apt purge nvidia-* 
sudo apt autoremove
sudo apt-get install nvidia-driver-510 -y
```

To verify version
```ls /usr/src | grep nvidia```

use  ```nvidia-smi``` to do the final check

If get ```Failed to initialize NVML: Driver/library version mismatch```, reboot and retry

If get ```NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver. Make sure that the latest NVIDIA driver is installed and running.```, the cause may be the card is too old. Latest driver does not support.

## 2. Install cuda tool kit
Follow the commands from [nvidia website](https://developer.nvidia.com/cuda-downloads)

Latest 11.6 version cude required driver version 510, so if your cards are old, dont install the latest 11.6 version.

```sudo apt-get install cuda=11.4.4-1``` works with driver 470
