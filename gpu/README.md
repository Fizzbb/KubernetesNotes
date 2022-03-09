# GPU node preparation

The latest cuda driver 510 only support those [cards](https://docs.nvidia.com/datacenter/tesla/tesla-release-notes-510-47-03/index.html), if you have older card, e.g., K80, install 470.

## 1. Install Nvidia Driver (must have)
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

## 3. Install Nvidia container runtime (for using GPUs inside container)
Nvidia contaienr runtime is a modified version of runc adding a custom pre-start hook to all containers. By using it, all the gpu devices are exposed to container.
Refer to Nvidia's [github](https://github.com/NVIDIA/nvidia-container-runtime) and [installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker) sites.
```
distribution=$(. /etc/os-release;echo $ID$VERSION_ID) \
   && curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add - \
   && curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
```
```
sudo apt-get update
sudo apt-get install -y nvidia-docker2
```
After install nvidia docker runtime, you can verify the gpu access from contaienrs with the following command. 

```docker run –rm –runtime=nvidia nvidia/cuda:11.0-base nvidia-smi```

We can also configure nvidia docker runtime as default (replace runc) by add the followings to the file /etc/docker/daemon.json and restart docker ```systemctl restart docker```
```
{
    "default-runtime":"nvidia",
    "runtimes": {
        "nvidia": {
            "path": "/usr/bin/nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
}
```
**WARNINFS: if use nvidia docker container and [Nvidia device plugin](https://github.com/NVIDIA/k8s-device-plugin#running-gpu-jobs) in the Kubernetes, and users dont set the GPU requirements in their yaml file, all the gpus will be exposed, which may cause resource management conflicts.**

## 4. Install cuDNN (for running deep learning workloads)
Install cuDNN on the host, if you want to directly run training on the host machine. Often time, we use containers with pre-built images which come with this library already.

The NVIDIA CUDA® Deep Neural Network library (cuDNN) is a GPU-accelerated library of primitives for deep neural networks. cuDNN provides highly tuned implementations for standard routines such as forward and backward convolution, pooling, normalization, and activation layers.

Deep learning researchers and framework developers worldwide rely on cuDNN for high-performance GPU acceleration. If not installed, errors like those ```Could not load dynamic library 'libcudnn.so.8'``` will show.

You can download the deb file from Nvidia [website](https://developer.nvidia.com/rdp/form/cudnn-download-survey), need to fill out a survey first.

Runtime library should be enough, for example 
```cuDNN Runtime Library for Ubuntu18.04 x86_64 (Deb)```.

Install is like this.

```sudo dpkg -i libcudnn8_8.2.4.15-1+cuda11.4_amd64.deb```
