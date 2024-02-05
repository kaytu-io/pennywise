# Usage Parameters

## AWS

#### azurerm_api_management
- self_hosted_gateway_count
- monthly_api_calls

#### azurerm_app_service_environment
- operating_system

#### azurerm_automation_account
- monthly_job_run_mins
- non_azure_config_node_count
- monthly_watcher_hrs

#### azurerm_automation_dsc_configuration
- non_azure_config_node_count

#### azurerm_automation_dsc_nodeconfiguration
- non_azure_config_node_count

#### azurerm_automation_job_schedule
- monthly_job_run_mins

#### azurerm_cdn_endpoint
- monthly_outbound_gb
- monthly_rules_engine_requests

#### azurerm_search_service
- monthly_images_extracted

#### azurerm_container_registry
- storage_gb
- monthly_build_vcpu_hrs

#### azurerm_cosmosdb_cassandra_keyspace
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_cassandra_table
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_gremlin_database
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb
 
#### azurerm_cosmosdb_gremlin_graph
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb
 
#### azurerm_cosmosdb_mongo_collection
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_mongo_database
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_sql_container
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_sql_database
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_cosmosdb_table
- monthly_serverless_request_units
- max_request_units_utilization_percentage
- monthly_analytical_storage_read_operations
- monthly_analytical_storage_write_operations
- storage_gb
- monthly_restored_data_gb

#### azurerm_dns_a_record
- monthly_queries

#### azurerm_dns_aaaa_record
- monthly_queries

#### azurerm_dns_caa_record
- monthly_queries

#### azurerm_dns_cname_record
- monthly_queries

#### azurerm_dns_mx_record
- monthly_queries

#### azurerm_dns_ns_record
- monthly_queries

#### azurerm_dns_ptr_record
- monthly_queries

#### azurerm_dns_srv_record
- monthly_queries

#### azurerm_dns_txt_record
- monthly_queries

#### azurerm_function_app
- monthly_executions
- execution_duration_ms
- memory_mb
- instances

#### azurerm_linux_function_app
- monthly_executions
- execution_duration_ms
- memory_mb
- instances

#### azurerm_windows_function_app
- monthly_executions
- execution_duration_ms
- memory_mb
- instances

#### azurerm_key_vault_certificate
- monthly_certificate_renewal_requests
- monthly_certificate_other_operations

#### azurerm_key_vault_key
- monthly_secrets_operations
- monthly_key_rotation_renewals
- monthly_protected_keys_operations
- hsm_protected_keys

#### azurerm_kubernetes_cluster
- nodes
- monthly_hrs
- monthly_data_processed_gb

#### azurerm_kubernetes_cluster_node_pool
- nodes
- monthly_hrs

#### azurerm_lb
- monthly_data_processed_gb

#### azurerm_mariadb_server
- additional_backup_storage_gb

#### azurerm_mssql_database
- extra_data_storage_gb
- monthly_vcore_hours
- long_term_retention_storage_gb
- backup_storage_gb

#### azurerm_mssql_managed_instance
- long_term_retention_storage_gb
- backup_storage_gb

#### azurerm_mysql_server
- additional_backup_storage_gb

#### azurerm_mysql_flexible_server
- additional_backup_storage_gb

#### azurerm_postgresql_flexible_server
- additional_backup_storage_gb

#### azurerm_postgresql_server
- additional_backup_storage_gb

#### azurerm_sql_database
- extra_data_storage_gb
- monthly_vcore_hours
- long_term_retention_storage_gb
- backup_storage_gb

#### azurerm_sql_managed_instance
- long_term_retention_storage_gb
- backup_storage_gb

#### azurerm_storage_account
- storage_gb
- monthly_iterative_read_operations
- monthly_read_operations
- monthly_iterative_write_operations
- monthly_write_operations
- monthly_list_and_create_container_operations
- monthly_other_operations
- monthly_data_retrieval_gb
- monthly_data_write_gb
- blob_index_tags
- data_at_rest_storage_gb
- snapshots_storage_gb
- metadata_at_rest_storage_gb
- early_deletion_gb

#### azure_storage_queue
- monthly_storage_gb
- monthly_class_1_operations
- monthly_class_2_operations
- monthly_geo_replication_data_transfer_gb

#### azure_storage_share
- storage_gb
- snapshots_storage_gb
- monthly_read_operations
- monthly_write_operations
- monthly_list_operations
- monthly_other_operations
- monthly_data_retrieval_gb
- metadata_at_rest_storage_gb

#### azurerm_virtual_machine
- monthly_os_disk_operations
- monthly_data_disk_operations
- monthly_hours

#### azurerm_windows_virtual_machine
- monthly_hours

#### azurerm_linux_virtual_machine
- monthly_hours

#### azurerm_managed_disk
- monthly_disk_operations

#### azurerm_image
- storage_gb

#### azurerm_virtual_machine_scale_set
- monthly_hours
- os_disk_monthly_operations
- data_disk_monthly_operations
- instances

#### azurerm_windows_virtual_machine_scale_set
- monthly_hours
- os_disk_monthly_operations

#### azurerm_linux_virtual_machine_scale_set
- monthly_hours
- os_disk_monthly_operations

#### azurerm_virtual_network_peering
- monthly_data_transfer_gb

#### azurerm_nat_gateway
- monthly_data_processed_gb

#### azurerm_private_endpoint
- monthly_inbound_data_processed_gb
- monthly_outbound_data_processed_gb

#### azurerm_virtual_network_gateway
- p2s_connection
- monthly_data_transfer_gb