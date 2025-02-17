import sys
from azure.ai.ml import MLClient
from azure.identity import DefaultAzureCredential
from azure.ai.ml.entities import Project
from dotenv import load_dotenv
import os

# Load environment variables from .env file
load_dotenv()

# Get CLI arguments
if len(sys.argv) < 3:
    print("Usage: python create_project.py <project_name> <hub_name>")
    sys.exit(1)

project_name = sys.argv[1]
hub_name = sys.argv[2]

# Get subscription ID and resource group from environment variables
subscription_id = os.getenv("SUBSCRIPTION_ID")
resource_group = os.getenv("RESOURCE_GROUP")

if not subscription_id or not resource_group:
    print("❌ SUBSCRIPTION_ID and RESOURCE_GROUP must be set in the .env file")
    sys.exit(1)

# Authenticate with Azure
credential = DefaultAzureCredential()
ml_client = MLClient(credential, subscription_id, resource_group)

# Define the project
hub_id = f"/subscriptions/{subscription_id}/resourceGroups/{resource_group}/providers/Microsoft.MachineLearningServices/workspaces/{hub_name}"
project = Project(name=project_name, hub_id=hub_id)

# Create the project
created_project = ml_client.workspaces.begin_create(workspace=project).result()
print(f"✅ Project '{created_project.name}' created successfully!")