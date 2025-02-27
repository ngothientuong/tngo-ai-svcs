<!DOCTYPE html>
<html>
<head>
<style>
  .center {
    text-align: center;
  }
  .pink-cursive {
    color: pink;
    font-family: "Brush Script MT", cursive;
    font-size: 72px;
    font-weight: bold;
  }
</style>
</head>
<body>

<div class="center">
  <span class="pink-cursive">Azure Cloud Setup</span>
</div>

</body>
</html>


# **🚀 Azure Cloud Setup from Scratch**
A fully automated CLI workflow to create an Azure subscription, resources, App Registrations, Service Principals, and assign permissions.

---


## **🟢 1️⃣ Create a New Azure Subscription (If Needed)**
    az account subscription create --display-name "My New Subscription"

#### **View All Subscriptions**
    az account list --query "[].{SubscriptionName:name, SubscriptionID:id}" -o table

- Use the **Subscription ID** for later commands.

---

## **🟢 2️⃣ Create a Resource Group**
    az group create --name myResourceGroup --location eastus

#### **View All Resource Groups**
    az group list --query "[].{Name:name, Location:location}" -o table

---

## **🟢 3️⃣ Create an App Registration & Its Service Principal**
### **(This will create the app + the SP at once)**
    az ad app create --display-name "MyApp"

#### **List All App Registrations (Applications)**
    az ad app list --query "[].{AppName:displayName, AppID:appId}" -o table

- Note the **Application (Client) ID**.

---

## **🟢 4️⃣ Create a Service Principal for the App**
    az ad sp create --id ${AZURE_SVC_PRINCIPAL_APP_ID}

#### **List All Service Principals (Enterprise Applications)**
    az ad sp list --query "[].{SPName:displayName, SPID:id, AppID:appId}" -o table

- Note the **Service Principal ID (SPID)**.

#### **Find the Service Principal for a Specific App**
    az ad sp show --id ${AZURE_SVC_PRINCIPAL_APP_ID} --query "{Name:displayName, SPID:id, AppID:appId}"

---

## **🟢 5️⃣ Assign Full API Permissions to the App**

Assign all required API permissions, including Microsoft Graph and Azure Service Management, with admin consent.

**Example**:
    `az ad app permission add --id ${AZURE_SVC_PRINCIPAL_APP_ID} --api <API_ID> --api-permissions <PERMISSION_ID>=Role`

---

**Step 1**: Assign Full Microsoft Graph API Permissions

    az ad app permission add --id ${AZURE_SVC_PRINCIPAL_APP_ID} --api 00000003-0000-0000-c000-000000000000 \
        --api-permissions 7ab1d382-f21e-4acd-a863-ba3e13f7da61=Role \
                         df021288-bdef-4463-88db-98f22de89214=Role \
                         19dbc75e-c2e2-444c-a770-ec69d8559fc7=Role \
                         06da0dbc-49e2-44d2-8312-53f166ab848a=Role \
                         e1fe6dd8-ba31-4d61-89e7-88639da4683d=Role \
                         741f803b-c850-494e-b5df-cde7c675a1ca=Role \
                         3afa6a7d-9b1a-42eb-948e-1650a849e176=Role \
                         3b5c71d1-d41e-4ccf-9952-7cd3f43c5b3e=Role \
                         e5e14a68-cf7c-4bbd-b7dd-d0513f6c04d0=Role



| **Permission** | **Permission ID** | **Description** |
|---------------|------------------|----------------|
| Directory.ReadWrite.All | `7ab1d382-f21e-4acd-a863-ba3e13f7da61` | Read and write all directory data |
| User.ReadWrite.All | `df021288-bdef-4463-88db-98f22de89214` | Read and write all users’ full profiles |
| RoleManagement.ReadWrite.Directory | `19dbc75e-c2e2-444c-a770-ec69d8559fc7` | Read and write directory roles |
| AppRoleAssignment.ReadWrite.All | `06da0dbc-49e2-44d2-8312-53f166ab848a` | Manage app role assignments |
| Application.ReadWrite.All | `e1fe6dd8-ba31-4d61-89e7-88639da4683d` | Read and write all applications |
| Group.ReadWrite.All | `741f803b-c850-494e-b5df-cde7c675a1ca` | Read and write all groups |
| Directory.Read.All	  | `3b5c71d1-d41e-4ccf-9952-7cd3f43c5b3e` |	Read directory data (needed for SPs, users, roles) |
| RoleManagement.Read.All | `e5e14a68-cf7c-4bbd-b7dd-d0513f6c04d0` | Read role assignments & Azure AD roles |

---

**Step 3**: Grant API Permissions to Make Changes Effective

    az ad app permission grant --id ${AZURE_SVC_PRINCIPAL_APP_ID} --api 00000003-0000-0000-c000-000000000000  \
      --scope Directory.ReadWrite.All User.ReadWrite.All RoleManagement.ReadWrite.Directory AppRoleAssignment.ReadWrite.All Application.ReadWrite.All Group.ReadWrite.All User.Read.All User.ReadWrite.All 06da0dbc-49e2-44d2-8312-53f166ab848a 3afa6a7d-9b1a-42eb-948e-1650a849e176  Directory.Read.All Directory.Read.All User.Read.All e1fe6dd8-ba31-4d61-89e7-88639da4683d RoleManagement.ReadWrite.All RoleManagement.Read.All

**Step 4** Grant permission for `Azure Service Management`

    az ad app permission add --id ${AZURE_SVC_PRINCIPAL_APP_ID} --api 797f4846-ba00-4fd7-ba43-dac1f8f63013 \
        --api-permissions 41094075-9dad-400e-a0bd-54e686782033=Role

    az ad app permission grant --id ${AZURE_SVC_PRINCIPAL_APP_ID} --api 797f4846-ba00-4fd7-ba43-dac1f8f63013 --scope user_impersonation


| **Permission** | **Permission ID** | **Description** |
|---------------|------------------|----------------|
| Access Azure Service Management | `41094075-9dad-400e-a0bd-54e686782033` | Manage Azure subscriptions and resources |

**Step 5**: Test Additional Graph API Calls

    # List all Groups
    az rest --method GET --uri "https://graph.microsoft.com/v1.0/groups" \
        --headers "Authorization=Bearer $(az account get-access-token --resource https://graph.microsoft.com/ --query accessToken -o tsv)"

    # List all Applications
    az rest --method GET --uri "https://graph.microsoft.com/v1.0/applications" \
        --headers "Authorization=Bearer $(az account get-access-token --resource https://graph.microsoft.com/ --query accessToken -o tsv)"

**Step 6**: Grant Admin Consent (Required for Some APIs)

    az ad app permission admin-consent --id ${AZURE_SVC_PRINCIPAL_APP_ID}

---

**Step 7**: List API Permissions for an App

    az ad app permission list --id ${AZURE_SVC_PRINCIPAL_APP_ID} -o table

---



## **🟢 6️⃣ Assign Roles to the Service Principal**
#### **Assign 'User Access Administrator' Role at Subscription Level**
    az role assignment create --assignee <SPID> --role "User Access Administrator" --scope /subscriptions/<SUBSCRIPTION_ID>

#### **Assign Additional Privileged Roles**
    az role assignment create --assignee <SPID> --role "Owner" --scope /subscriptions/<SUBSCRIPTION_ID>
    az role assignment create --assignee <SPID> --role "Storage Blob Data Owner" --scope /subscriptions/<SUBSCRIPTION_ID>

#### **View All Role Assignments for a Specific SP**
    az role assignment list --assignee <SPID> --all -o table

---

## **🟢 7️⃣ Create a Managed Identity for a Resource**
    az identity create --name MyManagedIdentity --resource-group myResourceGroup

#### **View All Managed Identities**
    az identity list --query "[].{Name:name, ID:id, PrincipalID:principalId, TenantID:tenantId}" -o table

- Note **Principal ID** for role assignments.

#### **Assign 'Managed Identity Operator' Role**
    az role assignment create --assignee <PRINCIPAL_ID> --role "Managed Identity Operator" --scope /subscriptions/<SUBSCRIPTION_ID>

---

## **🟢 8️⃣ List All Relevant Identities & Permissions**
#### **View All App Registrations**
    az ad app list --query "[].{AppName:displayName, AppID:appId}" -o table

#### **View All Service Principals**
    az ad sp list --query "[].{SPName:displayName, SPID:id, AppID:appId}" -o table

#### **View All Users in Azure AD**
    az ad user list --query "[].{User:displayName, Email:mail, ID:id}" -o table

#### **View All Managed Identities**
    az identity list --query "[].{Name:name, ID:id, PrincipalID:principalId, TenantID:tenantId}" -o table

#### **View All Role Assignments Across Subscription**
    az role assignment list --all --query "[].{Principal:principalName, Role:roleDefinitionName, Scope:scope, Type:principalType}" -o table

#### **View Role Assignments for a Specific Identity**
    az role assignment list --all --assignee "<IDENTITY_OBJECT_ID (or PRINCIPAL ID)>" --query "[].{Role:roleDefinitionName, Scope:scope}" -o table

---

## **✅ Summary of CLI Commands**
| **Command** | **Description** |
|-------------|----------------|
| `    az account list` | List all **Azure Subscriptions** |
| `    az group list` | List all **Resource Groups** |
| `    az ad app list` | List all **App Registrations** (Applications) |
| `    az ad sp list` | List all **Service Principals** |
| `    az role assignment list  --all --assignee <SP_OBJECT_ID> --output table` | List **role assignments** for a specific Service Principal |
| `    az role assignment list --all` | List **all role assignments** in the subscription |
| `    az ad user list` | List **all users in Azure AD** |
| `    az identity list --output table` | List **all Managed Identities** |
| `    az role assignment list --all --query "[].{ObjectID_OR_PrincipalID:principalId, PrincipalType:principalType, Role:roleDefinitionName, Scope:scope}" -o table` | List **all role assignments** for all identities |
| `    az role assignment list --all --query "[].{ObjectID_OR_PrincipalID:principalId, PrincipalType:principalType, Role:roleDefinitionName, Scope:scope}" -o table --assignee <IDENTITY_OBJECT_ID>` | View **role assignments** for a specific **User/SP/Identity** |
| `    az ad sp show --id $(az account show --query "user.name" -o tsv) -o table` | View current logged in **Service Principal (SP)** |
| `    az role assignment list --all --assignee $(az ad sp list --filter "appId eq '$(az account show --query user.name -o tsv)'" --query "[].id" -o tsv) --query "[].{ObjectID_OR_PrincipalID:principalId, PrincipalType:principalType, Role:roleDefinitionName, Scope:scope}" -o table` | Current role assignment list for logged in **Service Principal (SP)** |
---
