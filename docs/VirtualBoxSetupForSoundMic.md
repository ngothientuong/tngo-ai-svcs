# VirtualBox Ubuntu 24.04.1 Setup Guide with Working Audio & Mic

## Step 1: Install VirtualBox & VirtualBox Extension Pack

1. Download and install [VirtualBox](https://www.virtualbox.org/).
2. Download and install [VirtualBox Extension Pack](https://www.virtualbox.org/wiki/Downloads).

## Step 2: Download Ubuntu 24.04.1 ISO

1. Go to [Ubuntu Official Website](https://ubuntu.com/download/desktop).
2. Download **Ubuntu 24.04.1 LTS Desktop ISO**.

## Step 3: Create a New Virtual Machine

1. Open VirtualBox and click **New**.
2. Name it **tuongubuntu**.
3. Set **Type**: Linux, **Version**: Ubuntu (64-bit).
4. Allocate at least **20GB Storage** (Dynamically Allocated).
5. Set **Base Memory**: 20588 MB.
6. Set **Processors**: 10.
7. Click **Create**.

## Step 4: Configure Virtual Machine Settings

1. Open **Settings**.
2. Under **System → Motherboard**, set:
   - Boot Order: Hard Disk, Optical.
   - **Acceleration**: Enable Nested Paging, PAE/NX, KVM Paravirtualization.
3. Under **Display**:
   - **Video Memory**: 256MB.
   - **Graphics Controller**: VMSVGA.
   - Enable **3D Acceleration**.
4. Under **Storage**:
   - Add **Ubuntu 24.04.1 ISO** to IDE Controller.
   - Add **VBoxGuestAdditions ISO**.
5. Under **Audio**:
   - **Host Driver**: Windows DirectSound.
   - **Controller**: Intel HD Audio.
6. Under **Network**:
   - **Adapter 1**: Bridged Adapter.
   - **Name**: Intel(R) I211 Gigabit Network Connection.
   - **Adapter Type**: Intel PRO/1000 MT Desktop (82540EM).
   - **Promiscuous Mode**: Allow All.
   - Enable **Cable Connected**.

## Step 5: Install Ubuntu 24.04.1

1. Start the Virtual Machine.
2. Select **Install Ubuntu**.
3. Follow the installation process:
   - Select **Minimal Installation**.
   - Enable **Download Updates** and **Install Third-Party Software**.
   - Choose **Erase Disk and Install Ubuntu**.
4. Restart when prompted.

## Step 6: Install VirtualBox Guest Additions

1. Open Terminal in Ubuntu (`Ctrl+Alt+T`).
2. Run:
   ```bash
   sudo apt update -y && sudo apt upgrade -y
   sudo apt install -y build-essential dkms linux-headers-$(uname -r)
   ```
3. Insert Guest Additions CD:
   - In VirtualBox, go to **Devices → Insert Guest Additions CD Image**.
4. Mount and install:
   ```bash
   sudo mount /dev/cdrom /mnt
   cd /mnt
   sudo ./VBoxLinuxAdditions.run
   ```
5. Reboot:
   ```bash
   sudo reboot
   ```

## Step 7: Install Audio & Video Prerequisites

1. Open Terminal and run:
   ```bash
   sudo apt-get install -y pulseaudio-utils pulseaudio vim
   sudo apt install ubuntu-restricted-extras gstreamer1.0-libav gstreamer1.0-plugins-bad at-spi2-core -y
   ```

## Step 8: Restart Audio Services After Every Login

1. **Every time you log in, manually run:**
   ```bash
   pulseaudio --kill
   pulseaudio --start
   sudo alsa force-reload
   pactl list sinks short
   systemctl --user restart pipewire pipewire-pulse
   ```
2. **Test microphone recording:**
   ```bash
   arecord -D pulse -f cd -d 5 test-mic.wav
   aplay test-mic.wav
   ```
   ✅ If you hear your voice, **audio & mic are working!** 🎤🔥

## Step 9: (Optional) Automate Audio Restart at Boot

If you want to **automate the audio fix** at every startup:

1. Create an audio restart script:
   ```bash
   sudo vim /usr/local/bin/audio-reload.sh
   ```
2. Paste this inside:
   ```bash
   #!/bin/bash
   pulseaudio --kill
   pulseaudio --start
   sudo alsa force-reload
   pactl list sinks short
   systemctl --user restart pipewire pipewire-pulse
   ```
3. Make it executable:
   ```bash
   sudo chmod +x /usr/local/bin/audio-reload.sh
   ```
4. Add it to cron so it runs at every boot:
   ```bash
   sudo crontab -e
   ```
   Add this line at the bottom:
   ```bash
   @reboot /usr/local/bin/audio-reload.sh
   ```
5. **Now reboot and test!**
   ```bash
   sudo reboot
   ```

## 🎉 Conclusion
✅ **Your VirtualBox Ubuntu 24.04.1 now has fully working audio & mic!** 🎤🔥
✅ **This guide ensures proper setup, audio packages, and manual steps for every reboot!** 🚀
✅ **If you followed Step 9, your audio will also auto-restart after every boot!**

🔥 **Enjoy your VirtualBox Ubuntu with working sound & mic!** 🔥