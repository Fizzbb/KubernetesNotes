# GPU node preparation

The latest cuda driver 510 only support those [cards](https://docs.nvidia.com/datacenter/tesla/tesla-release-notes-510-47-03/index.html), if you have older card, e.g., K80, install 470.

## 1. Install Nvidia Driver
remove existing one and install the latest
```
sudo apt clean
sudo apt update
sudo apt purge nvidia-* 
sudo apt autoremove
sudo apt-get install nvidia-driver-510 -y
```

To verify driver version
```ls /usr/src | grep nvidia```

use  ```nvidia-smi``` to do the final check

If get Error 1: ```Failed to initialize NVML: Driver/library version mismatch```, reboot and retry

If get Error 2: ```NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver. Make sure that the latest NVIDIA driver is installed and running.```, the cause may be the card is too old. Latest driver does not support.

## 2. Install CUDA tool kit (for build/profile... cuda program)
Follow the commands from [nvidia website](https://developer.nvidia.com/cuda-downloads)

However, if your cards are Kepler or older, which are not supported by driver verision 510, don't directly install cuda like the following command mentioned in above website.

```sudo apt-get install cuda=11.4.4-1```

The cuda package will automatically remove older driver, and install 510, which will lead to Error 2 above. 

Instead of installing cuda, we can install cuda-toolkit, which won't automatically update driver version. Refer to this [page](https://forums.developer.nvidia.com/t/cuda-11-4-installer-wants-to-install-nvidia-driver-version-incompatible-with-tesla-k40m/192879) for more info. 

```sudo apt install cuda-toolkit-11-4```

Lastly, add cuda path to your path

```export PATH=$PATH:/usr/local/cuda/bin```

To verify installed cuda version
```nvcc --version```
