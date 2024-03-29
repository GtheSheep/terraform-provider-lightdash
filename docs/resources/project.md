---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lightdash_project Resource - terraform-provider-lightdash"
subcategory: ""
description: |-
  
---

# lightdash_project (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dbt_connection_repository` (String) Repository name in <org>/<repo> format
- `name` (String) Project name
- `organization_uuid` (String, Sensitive) UUID of the organization to create the project in
- `type` (String) Type of project to create, either DEFAULT or DEVELOPMENT
- `warehouse_connection_account` (String) Account identifier, including region/ cloud path
- `warehouse_connection_database` (String) Database to connect to
- `warehouse_connection_role` (String) Role to connect to the warehouse with
- `warehouse_connection_warehouse` (String) Warehouse to use

### Optional

- `dbt_connection_branch` (String) Branch to use, default 'main'
- `dbt_connection_host_domain` (String) Host domain of the repo, default 'github.com'
- `dbt_connection_project_sub_path` (String) Sub path to find the project in the repo, default '/'
- `dbt_connection_type` (String) dbt project connection type, currently only support 'github', which is the default
- `warehouse_connection_client_session_keep_alive` (Boolean) Client session keep alive param, default `false`
- `warehouse_connection_schema` (String) Schema to connect to, default 'PUBLIC'
- `warehouse_connection_threads` (Number) Number of threads to use, default `1`
- `warehouse_connection_type` (String) Type of warehouse to connect to, currently only 'snowflake', as a default

### Read-Only

- `id` (String) The ID of this resource.
