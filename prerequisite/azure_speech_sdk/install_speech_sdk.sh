#!/bin/bash

sudo apt-get update
sudo apt-get install build-essential ca-certificates libasound2-dev libssl-dev wget -y
# Live Streaming
sudo apt install ffmpeg -y
# Ensure to install `python ahead of time!`
pip install yt-dlp

# Load environment variables from .env file
ENV_FILE=$(dirname $(dirname $(dirname $(realpath $0))))/.env

if [ ! -f "$ENV_FILE" ]; then
  echo ".env file not found at $ENV_FILE"
  exit 1
fi

# Source the environment variables from the .env file
set -o allexport
source $ENV_FILE
set +o allexport

# Check if SPEECH_LINUX_SDK_VERSION is set
if [ -z "$SPEECH_LINUX_SDK_VERSION" ]; then
  echo "SPEECH_LINUX_SDK_VERSION is not set in the environment variables"
  exit 1
fi

# Define the installation directory
INSTALL_DIR=$(dirname $(dirname $(dirname $(realpath $0))))/sdk/azure_speech_sdk

# Create the installation directory
mkdir -p $INSTALL_DIR

# Download the specific version of the Azure Speech SDK
echo "Downloading Azure Speech SDK version $SPEECH_LINUX_SDK_VERSION..."
wget https://csspeechstorage.blob.core.windows.net/drop/$SPEECH_LINUX_SDK_VERSION/SpeechSDK-Linux-$SPEECH_LINUX_SDK_VERSION.tar.gz -O SpeechSDK-Linux-$SPEECH_LINUX_SDK_VERSION.tar.gz

# Extract the SDK
echo "Extracting Azure Speech SDK version $SPEECH_LINUX_SDK_VERSION..."
tar -xzf SpeechSDK-Linux-$SPEECH_LINUX_SDK_VERSION.tar.gz -C $INSTALL_DIR
rm SpeechSDK-Linux-$SPEECH_LINUX_SDK_VERSION.tar.gz

# Verify extraction
echo "Contents of $INSTALL_DIR:"
ls -l $INSTALL_DIR

# Set environment variables
export SPEECHSDK_ROOT=$INSTALL_DIR/SpeechSDK-Linux-$SPEECH_LINUX_SDK_VERSION
export LD_LIBRARY_PATH=$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH
export C_INCLUDE_PATH=$SPEECHSDK_ROOT/include/c_api:$C_INCLUDE_PATH

# Make the environment variables persistent across reboots
echo "export SPEECHSDK_ROOT=$SPEECHSDK_ROOT" >> ~/.bashrc
echo "export LD_LIBRARY_PATH=$SPEECHSDK_ROOT/lib/x64:\$LD_LIBRARY_PATH" >> ~/.bashrc
echo "export C_INCLUDE_PATH=$SPEECHSDK_ROOT/include/c_api:\$C_INCLUDE_PATH" >> ~/.bashrc
echo "export CGO_CFLAGS=\"-I$SPEECHSDK_ROOT/include/c_api\"" >> ~/.bashrc
echo "export CGO_LDFLAGS=\"-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core\"" >> ~/.bashrc
source ~/.bashrc

echo "Azure Speech SDK version $SPEECH_LINUX_SDK_VERSION installed successfully in $INSTALL_DIR."