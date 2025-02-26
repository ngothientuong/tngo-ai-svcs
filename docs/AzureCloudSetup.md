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

## **🔹 CLI Commands for Viewing All Apps, Service Principals, and Role Assignments**

---

## **1️⃣ View All App Registrations (Applications)**
    az ad app list --query "[].{AppName:displayName, AppID:appId}" -o table

- This lists all **App Registrations** with their **Name** and **Application (Client) ID**.
- Use `-o json` if you need detailed output.

---

## **2️⃣ View All Service Principals (Enterprise Applications)**
    az ad sp list --query "[].{SPName:displayName, SPID:id, AppID:appId}" -o table

- This lists all **Service Principals** along with their **Object ID (SPID)** and **App ID**.

#### **Find the Service Principal for a Specific App:**
    az ad sp show --id <APP_ID>

- Replace `<APP_ID>` with the **Application (Client) ID**.

---

## **3️⃣ View All Role Assignments for a Specific App or Service Principal**
    az role assignment list --assignee <SP_OBJECT_ID> --query "[].{Role:roleDefinitionName, Scope:scope, AssignedTo:principalName}" -o table

- Replace `<SP_OBJECT_ID>` with the **Service Principal Object ID**.

---

## **4️⃣ View All Role Assignments Across Subscription**
    az role assignment list --query "[].{Role:roleDefinitionName, Scope:scope, AssignedTo:principalName}" -o table

- This gives a full list of **who has what roles at what scope**.

---

## **5️⃣ View All Identities (Users, Service Principals, and Managed Identities)**
    az ad sp list --query "[].{Name:displayName, Type:servicePrincipalType, ID:id}" -o table

- This shows **Service Principals, Managed Identities**, and their **Object IDs**.

#### **View All Users in Azure AD**
    az ad user list --query "[].{User:displayName, Email:mail, ID:id}" -o table

- This lists all **Azure AD users**.

---

## **6️⃣ View All Role Assignments for All Identities**
    az role assignment list --all --query "[].{Principal:principalName, Role:roleDefinitionName, Scope:scope, Type:principalType}" -o table

- This lists all **Users, Service Principals, and Managed Identities** with their **assigned roles**.

---

## **7️⃣ View Role Assignments for a Specific Identity (User, SP, or Managed Identity)**
    az role assignment list --assignee "<IDENTITY_OBJECT_ID>" --query "[].{Role:roleDefinitionName, Scope:scope}" -o table

- Replace `<IDENTITY_OBJECT_ID>` with the **User ID, Service Principal ID, or Managed Identity ID**.

---

## **8️⃣ Get a Specific Role Assignment for Debugging**
    az role assignment list --assignee <SP_OBJECT_ID> --role "User Access Administrator" -o json

- This checks if a Service Principal has **"User Access Administrator"** role.

---

## **✅ Summary of CLI Commands**
| **Command** | **Description** |
|-------------|----------------|
| `    az ad app list` | List all **App Registrations** (Applications) |
| `    az ad sp list` | List all **Service Principals** |
| `    az role assignment list --assignee <SP_OBJECT_ID>` | List **role assignments** for a specific Service Principal |
| `    az role assignment list` | List **all role assignments** in the subscription |
| `    az ad user list` | List **all users in Azure AD** |
| `    az ad sp list --query "[].{Name:displayName, Type:servicePrincipalType, ID:id}" -o table` | List **all identities** (SPs, Managed Identities) |
| `    az role assignment list --all --query "[].{Principal:principalName, Role:roleDefinitionName, Scope:scope, Type:principalType}" -o table` | List **all role assignments** for all identities |
| `    az role assignment list --assignee <IDENTITY_OBJECT_ID>` | View **role assignments** for a specific **User/SP/Identity** |

