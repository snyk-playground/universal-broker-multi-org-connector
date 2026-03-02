# Universal Broker Multi-Org Connector

**Automate Snyk organization management and integrate a single Universal Broker connection across multiple organizations.**

![Snyk OSS Example](https://raw.githubusercontent.com/snyk-labs/oss-images/main/oss-example.jpg)

This repo provides:

- Extraction of Snyk groups
- Extraction of Snyk organizations by group(s)
- Bulk integration of a Universal Broker connection to multiple organizations
- Rollback of bulk integrations of a Universal Broker connection if needed

## Requirements

- Snyk account with API access
- **Environment variables only:**
  - `SNYK_API_TOKEN` (required)
  - `SNYK_TENANT_ID` (required)

## Installation

### Downloading a Release from GitHub

Visit the [Releases](https://github.com/snyk-playground/universal-broker-multi-org-connector/releases) page for the GitHub project, and find the appropriate archive for your
operating system and architecture. Download the archive from your browser to your home directory.

Move the `broker-moc` binary to somewhere in your path. For example, on GNU/Linux and macOS:

```bash
mv broker-moc /usr/local/bin
```

Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add
`broker-moc` to their PATH.

## Usage

### 1. Export credentials

```bash
export SNYK_API_TOKEN=your_api_token_here
export SNYK_TENANT_ID=your_tenant_id_here
```

### 2. List organizations

The `org list` command lists all organizations accessible to your user. You can save the output
to a file for use in the integration step.

> Tip: Run `broker-moc org list --help` to see all available flags.

To view all accessible organizations:

```bash
$ broker-moc org list
```

<img alt="list accessible organizations" src=".github/cmd_org_list.gif" />

To filter organizations by specific groups, provide the `--group-id` flag (repeatable):

```bash
$ broker-moc org list --group-id <first-group-id> --group-id <second-group-id>
```

To save the organization list for use in the next step, add `--output` and `--format` flags:

```bash
$ broker-moc org list --group-id <first-group-id> --format yaml --output snyk_orgs.yaml
```

<details>
<summary>Example: snyk_orgs.yaml</summary>

```yaml
organizations:
  - id: 3d35a9a4-0948-4089-82b0-ff16140c1615
    name: demo-group-api_1
    slug: demo-group-api_1
    group_id: 300b8fb6-01a9-422a-8b06-044108d9ef42
    tenant_id: 2b2b2388-cbb1-4f55-964a-982e8b381530
    created_at: 2026-03-02T00:57:14Z
  - id: ca72fa3b-72a8-4b39-ace1-ea8456fb340f
    name: demo-group-api_2
    slug: demo-group-api_2
    group_id: 300b8fb6-01a9-422a-8b06-044108d9ef42
    tenant_id: 2b2b2388-cbb1-4f55-964a-982e8b381530
    created_at: 2026-03-02T00:57:24Z
  # commented out = skipped during integration
  #- id: 6e12cf10-b18c-44c8-bec1-ce42a131cff7
  #  name: demo-group-api_3
  #  slug: demo-group-api_3
  #  group_id: 300b8fb6-01a9-422a-8b06-044108d9ef42
  #  tenant_id: 2b2b2388-cbb1-4f55-964a-982e8b381530
  #  created_at: 2026-03-02T00:58:02Z
```

</details>

### 3. Connect a Broker connection to multiple organizations

Use the file generated in step 2 as input. The `--connection-type` value should match
the type configured in your Universal Broker Connection (e.g. github, github-enterprise, gitlab).

> Tip: Run `broker-moc connection integrate ... --dry-run` to print organizations without performing integration.

```bash
$ broker-moc connection integrate <your-connection-id> \
  --connection-type <your-connection-type> \
  --input snyk_orgs.yaml \
  --output snyk_connected_integrations.yaml \
  --format yaml
...
=> Starting integration process...
   organizations:   3
   tenant ID:       <your-tenant-id>
   connection ID:   <your-connection-id>
   connection type: <your-connection-type>
------------------------------------------------------------
=> Summary
   successful integrations:  3
   failed integrations:      0
------------------------------------------------------------
```

<img alt="integrate connection to multiple organizations" src=".github/cmd_connection_integrate.gif" />

<details>
<summary>Example: snyk_connected_integrations.yaml</summary>

```yaml
integrations:
  - id: 86511127-07b5-4b10-883c-368cb04292a0
    connection_id: 8b424fbf-b7a5-4ebd-b6e3-9bb056781956
    connection_type: gitlab
    organization_id: 3d35a9a4-0948-4089-82b0-ff16140c1615
    organization_name: demo-group-api_1
    tenant_id: 2b2b2388-cbb1-4f55-964a-982e8b381530
    status: connected
  - id: f4589c6e-b374-48d4-b501-e7966bfdead8
    connection_id: 8b424fbf-b7a5-4ebd-b6e3-9bb056781956
    connection_type: gitlab
    organization_id: ca72fa3b-72a8-4b39-ace1-ea8456fb340f
    organization_name: demo-group-api_2
    tenant_id: 2b2b2388-cbb1-4f55-964a-982e8b381530
    status: connected
```

</details>

### 4. Disconnect a Broker connection from organizations

To disconnect organizations that were connected in the previous step, reuse the output file from step 3 (e.g.
`snyk_connected_integrations.yaml`). Review the file and remove any organizations you want to keep connected, then run:

> Tip: Run `broker-moc connection disconnect ... --dry-run` to print organizations without performing disconnect.

```bash
$ broker-moc connection disconnect <your-connection-id> \
  --input=snyk_connected_integrations.yaml \
  --output=snyk_disconnected_integrations.yaml \
  --format yaml
...
=> Starting disconnect process...
   integrations:  3
   connection ID: <your-connection-id>
------------------------------------------------------------
=> Summary
   successful disconnections: 3
   failed disconnections:     0
------------------------------------------------------------
```

<img alt="disconnect connection from multiple integrations" src=".github/cmd_connection_disconnect.gif" />

## License

MIT
