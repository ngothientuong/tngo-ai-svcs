# PROJECT TO-DOs

## API Outstanding items
### Text To Speech - Custom Neural
- [Text To Speech - Custom Neural API](https://learn.microsoft.com/en-gb/rest/api/aiservices/speechapi/operation-groups?view=rest-aiservices-speechapi-2024-02-01-preview)
- Base Models
- Consents
- Endpoints
- Models
- Operations
- Personal Voices
- Training Sets
- **Other Requirements**:
  - Submit consents: `https://speech.microsoft.com/portal` -> `Custom Voice` -> `Access requirement` -> `Apply for limited access`
  - `Storage-sas`: Grant limited access to Azure Storage resources using shared access signatures (SAS)

### Text To Speech - Batch synthesis API
- [Batch synthesis API](https://learn.microsoft.com/en-us/azure/ai-services/speech-service/batch-synthesis)
- Upload large training set in batch via asynchronous fashion

### Storage-sas
- [How-to-Storage-Sas](https://learn.microsoft.com/en-us/azure/storage/common/storage-sas-overview)
- Instead of providing public URI for a audio, video files, you can store resources in Azure Storage and provide acccess via `shared access signatures` (`SAS`)
- A shared access signature (SAS) provides secure delegated access to resources in your storage account. With a SAS, you have granular control over how a client can access your data.


