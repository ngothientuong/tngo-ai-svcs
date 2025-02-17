# tngo-ai-svcs
Go repo for building AI services - Training repo

## Setup Instructions

1. Copy the provided `.env` file to your project directory and rename it to `.env`. Fill in the values for the variables as needed.

```sh
cp .env.example .env
```

2. **Important:** Do not commit the `.env` file to the repository in your branch. Add the `.env` file to your `.gitignore` to ensure it is not tracked by Git.

```sh
echo ".env" >> .gitignore
```

## Install Go

Follow these steps to install Go on your machine:

1. **Download the Go binary**

Visit the [official Go download page](https://golang.org/dl/) and download the appropriate binary for your operating system.

2. **Install Go**

Follow the installation instructions for your operating system:

- **Linux**

```sh
tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
```

- **macOS**

```sh
sudo tar -C /usr/local -xzf go1.23.1.darwin-amd64.tar.gz
```

- **Windows**

Run the MSI installer and follow the prompts.

3. **Set up Go environment variables**

Add the following lines to your `.bashrc`, `.zshrc`, or `.profile` file:

```sh
export PATH=$PATH:/usr/local/go/bin
```

Reload your shell configuration:

```sh
source ~/.bashrc
```

4. **Verify the installation**

Ensure that Go is installed correctly by running the following command:

```sh
go version
```

## Install Azure Speech SDK

Follow these steps to install the Azure Speech SDK:

1. **Run the installation script for Speech SDK**

Navigate to the `prerequisite/azure_speech_sdk` directory and run the `install_speech_sdk.sh` script where go mod is with `require github.com/Microsoft/cognitive-services-speech-sdk-go v1.33.0`:

```sh
cd prerequisite/azure_speech_sdk
./install_speech_sdk.sh
source ~/.bashrc

cd prerequisite/azure_ai
./install_python_ai.sh
```

2. **Install yt-dlp python package**
```python
pip install yt-dlp
```

This script will download and install the Azure Speech SDK and set the necessary environment variables.

## Azure Resources Setup

Follow these steps to create the necessary Azure resources and configure the service principal:

1. **Create a Resource Group**

```sh
az group create --name <resource-group-name> --location <location>
```

2. **Create an Azure Cognitive Services Account**

```sh
az cognitiveservices account create \
    --name <cognitive-services-account-name> \
    --resource-group <resource-group-name> \
    --kind CognitiveServices \
    --sku S0 \
    --location <location> \
    --yes
```

3. **Create a Service Principal**

```sh
az ad sp create-for-rbac --name <service-principal-name> --role Contributor \
    --scopes /subscriptions/<subscription-id>/resourceGroups/<resource-group-name> \
    --sdk-auth
```

This command will output a JSON object with the service principal credentials. Save this JSON object securely as it contains sensitive information.

4. **Set Permissions for the Service Principal**

Assign the service principal to the resource group with the necessary permissions:

```sh
az role assignment create --assignee <service-principal-id> \
    --role Contributor \
    --resource-group <resource-group-name>
```

5. **Configure Environment Variables**

Add the following environment variables to your `.env` file using the values from the previous steps:

```properties
AZURE_SVC_PRINCIPAL_APP_ID=<service-principal-app-id>
AZURE_SVC_PRINCIPAL_APP_PASSWORD=<service-principal-password>
AZURE_SVC_PRINCIPAL_TENANT_ID=<tenant-id>
AZURE_AI_MULTISERVICE_ENDPOINT=https://<cognitive-services-account-name>.cognitiveservices.azure.com
AZURE_AI_MULTISERVICE_KEY=<cognitive-services-account-key>
```

Replace the placeholders with the actual values obtained from the Azure portal and the service principal creation output.

6. **Verify the Setup**

Ensure that the environment variables are correctly set by running the following command:

```sh
source .env
```

You can now proceed with running your application using the configured Azure resources and service principal.

## Azure Cloud Permission

### Assign Managed Identities to a `source` resource
- When creating a resource such as AI Services, AI Foundry, Kubernetes..., you're often asked to create a system managed identity or user assign existing user managed identity to the `destination` resource.
- What it means is that those resources can utilize those managed identities to access OTHER Azure resources such as Storage Account, Container Registry...

### RBAC - Allow `source` Azure resource or `nonInteractive cli service principal (app)` to access `destination` resource via Managed identities
- At the `destination` resource such as Storage Account, Container Registry, Key Vault, you'll have to perform role assignment with the below rule:
  - `Role` relating to the resource such as `Storage Account Contributor` role for `Storage Account`, `Key Vault Administrator` role in `Key Vault`
  - `Who` is given the above `Role`:
    - `Source` resource to this `destination` resource: the user managed identities or system managed identities. Meaning those source resource assigned with these managed identities can obtain the role associating with this dest resource
    - `NonInteractive cli`: service principal (app)