# Workload Identity Federation: Cross-Organization Access

## Current Setup

Your deployment uses Workload Identity Federation (WIF) to authenticate GitHub Actions with GCP:

- **GCP Project**: `slap-ai-481400`
- **Service Account**: `ai-study-uploader@slap-ai-481400.iam.gserviceaccount.com`
- **Workload Identity Pool**: `github-deploy-pool`
- **Provider**: `github-provider`

## The Cross-Organization Problem

You mentioned changing the repository reference from your personal repo to `slap-events` org. This raises a valid concern about whether the workload identity provider can authenticate across different GitHub organizations.

### What WIF Actually Does

Workload Identity Federation allows GitHub Actions to exchange a GitHub OIDC token for a GCP access token **without storing long-lived credentials**. The key points:

1. **GitHub OIDC Token** contains claims about:
    - Repository owner/org
    - Repository name
    - Branch/ref
    - Actor (user who triggered the workflow)
    - Commit SHA

2. **WIF Provider** validates the OIDC token and maps it to the service account

3. **Attribute Mapping** in the provider determines which GitHub identities can assume the service account role

### Cross-Org Access: The Good News

**WIF can authenticate across organizations**, but it depends on your attribute mapping configuration.

#### Typical Attribute Mapping

```json
{
    "principalSetSelection": "expression",
    "expression": "google.subject == 'repo:OWNER/REPO:ref:refs/heads/BRANCH'"
}
```

Or more permissively:

```json
{
    "principalSetSelection": "expression",
    "expression": "google.subject.startsWith('repo:slap-events/')"
}
```

## Checking Your Current Configuration

### 1. View Your Workload Identity Provider

```bash
gcloud iam workload-identity-pools providers describe github-provider \
  --project=slap-ai-481400 \
  --location=global \
  --workload-identity-pool=github-deploy-pool \
  --format=json
```

Look for the `attributeMapping` and `attributeCondition` fields.

### 2. Check Service Account Bindings

```bash
gcloud iam service-accounts get-iam-policy \
  ai-study-uploader@slap-ai-481400.iam.gserviceaccount.com \
  --project=slap-ai-481400
```

Look for `roles/iam.workloadIdentityUser` bindings.

## Scenarios and Solutions

### Scenario 1: Repository in Same Organization (Current)

**Configuration**: Attribute mapping allows `repo:YOUR_ORG/learning-core-api`

**Status**: ✅ Works - WIF authenticates the workflow

**No action needed** - your current setup handles this.

---

### Scenario 2: Repository in Different Organization (`slap-events`)

**Problem**: If the WIF provider is configured to only allow your personal org, it will reject tokens from `slap-events` org workflows.

**Solution A: Update Attribute Mapping (Recommended)**

Update the workload identity provider to accept repositories from both organizations:

```bash
gcloud iam workload-identity-pools providers update-oidc github-provider \
  --project=slap-ai-481400 \
  --location=global \
  --workload-identity-pool=github-deploy-pool \
  --attribute-mapping="google.subject=assertion.sub,google.aud=assertion.aud" \
  --attribute-condition="assertion.repository_owner in ['YOUR_ORG', 'slap-events'] && assertion.repository == 'learning-core-api'"
```

**Solution B: Use Broader Mapping**

Allow any repository from specific organizations:

```bash
gcloud iam workload-identity-pools providers update-oidc github-provider \
  --project=slap-ai-481400 \
  --location=global \
  --workload-identity-pool=github-deploy-pool \
  --attribute-mapping="google.subject=assertion.sub,google.aud=assertion.aud" \
  --attribute-condition="assertion.repository_owner in ['YOUR_ORG', 'slap-events']"
```

---

### Scenario 3: Repository Access Across Different GCP Projects

**Problem**: If you need to deploy to a different GCP project (e.g., `slap-events` project), the service account in `slap-ai-481400` won't have permissions.

**Solution**:

1. **Option A**: Create a separate service account in the target project

    ```bash
    # In slap-events project
    gcloud iam service-accounts create github-deployer \
      --project=slap-events-project

    # Grant it necessary roles
    gcloud projects add-iam-policy-binding slap-events-project \
      --member="serviceAccount:github-deployer@slap-events-project.iam.gserviceaccount.com" \
      --role="roles/run.admin"
    ```

2. **Option B**: Grant cross-project permissions to existing service account
    ```bash
    # In target project, grant permissions to the slap-ai service account
    gcloud projects add-iam-policy-binding slap-events-project \
      --member="serviceAccount:ai-study-uploader@slap-ai-481400.iam.gserviceaccount.com" \
      --role="roles/run.admin"
    ```

---

## Verification Steps

### Step 1: Test WIF Authentication Locally

```bash
# Get the OIDC token from GitHub (only works in GitHub Actions)
# In your workflow, add a debug step:

- name: Debug WIF Token
  run: |
    TOKEN=$(curl -s -H "Authorization: bearer $ACTIONS_ID_TOKEN_REQUEST_TOKEN" \
      "$ACTIONS_ID_TOKEN_REQUEST_URL" | jq -r '.token')
    echo "Token claims:"
    echo "$TOKEN" | cut -d. -f2 | base64 -d | jq .
```

### Step 2: Verify Attribute Mapping

After running the debug step, check that the token contains:

- `repository_owner`: Should be your org name
- `repository`: Should be `learning-core-api`
- `ref`: Should be `refs/heads/main`

### Step 3: Check Service Account Permissions

```bash
# List all roles for the service account
gcloud projects get-iam-policy slap-ai-481400 \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:ai-study-uploader@*" \
  --format=table
```

Required roles:

- `roles/iam.workloadIdentityUser` - for WIF
- `roles/secretmanager.secretAccessor` - to read secrets
- `roles/artifactregistry.writer` - to push images
- `roles/run.admin` - to deploy to Cloud Run

---

## Current Deployment Workflow: What Actually Happens

1. **GitHub Actions runs** in your repository (any org)
2. **GitHub generates OIDC token** with claims about the repository
3. **Workflow calls** `google-github-actions/auth@v2` with WIF provider details
4. **Auth action exchanges** GitHub OIDC token for GCP access token
5. **GCP validates** the token against the WIF provider's attribute mapping
6. **If valid**, service account permissions are granted to the workflow
7. **Workflow can now**:
    - Access Secret Manager secrets
    - Push to Artifact Registry
    - Deploy to Cloud Run

---

## Recommended Configuration for Your Setup

Based on your needs (supporting both personal and `slap-events` org repos):

```bash
gcloud iam workload-identity-pools providers update-oidc github-provider \
  --project=slap-ai-481400 \
  --location=global \
  --workload-identity-pool=github-deploy-pool \
  --attribute-mapping="google.subject=assertion.sub,google.aud=assertion.aud,attribute.repository_owner=assertion.repository_owner" \
  --attribute-condition="assertion.repository_owner in ['YOUR_PERSONAL_ORG', 'slap-events'] && assertion.repository == 'learning-core-api'"
```

Then update the service account binding:

```bash
gcloud iam service-accounts add-iam-policy-binding \
  ai-study-uploader@slap-ai-481400.iam.gserviceaccount.com \
  --project=slap-ai-481400 \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/578398115938/locations/global/workloadIdentityPools/github-deploy-pool/attribute.repository_owner/YOUR_PERSONAL_ORG" \
  --member="principalSet://iam.googleapis.com/projects/578398115938/locations/global/workloadIdentityPools/github-deploy-pool/attribute.repository_owner/slap-events"
```

---

## Troubleshooting Checklist

- [ ] WIF provider attribute mapping includes both organizations
- [ ] Service account has `roles/iam.workloadIdentityUser` role
- [ ] Service account has `roles/secretmanager.secretAccessor` for all required secrets
- [ ] Service account has `roles/artifactregistry.writer` for the artifact registry
- [ ] Service account has `roles/run.admin` for Cloud Run deployment
- [ ] Secrets exist in Secret Manager with correct names
- [ ] GitHub repository is in an allowed organization (per attribute mapping)
- [ ] Workflow file syntax is valid YAML

---

## Summary

**Your current setup will work across organizations** as long as:

1. The WIF provider's attribute mapping includes the `slap-events` organization
2. The service account has the necessary IAM roles
3. All required secrets exist in Secret Manager

If you're getting authentication errors when deploying from `slap-events`, the first step is to verify the attribute mapping configuration using the commands above.
