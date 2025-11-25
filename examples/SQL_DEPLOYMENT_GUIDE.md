# SQL Deployment Examples

This directory contains examples for creating SQL deployments in Ververica Platform (VVP). There are two main approaches to running SQL deployments:

1. **Using a Deployment Target** - Dedicated resources for the deployment
2. **Using a Session Cluster** - Shared resources with other SQL jobs

## Example Files

### `sqldeployment.yaml`
Contains three examples:
- **Example 1**: SQL deployment with deployment target
- **Example 2**: Simple SQL deployment with session cluster
- **Example 3**: Complex SQL deployment with Kafka source and aggregations

### `test-sql-session.yaml`
A minimal test example for SQL deployments using a session cluster.

## Prerequisites

### For Session Cluster Deployments
1. First, create a session cluster:
   ```bash
   vvp2 sessioncluster create -n default -f examples/sessioncluster.yaml
   ```

2. Verify the session cluster is running:
   ```bash
   vvp2 sessioncluster list -n default
   ```

### For Deployment Target Deployments
1. Ensure you have a deployment target configured:
   ```bash
   vvp2 deployment-target list -n default
   ```

## Usage Examples

### 1. Create SQL Deployment with Session Cluster

```bash
# Using the simple example
vvp2 deployment create -n default -f examples/test-sql-session.yaml

# Using the analytics example
vvp2 deployment create -n default -f examples/sqldeployment.yaml
```

### 2. Create SQL Deployment with Deployment Target

```bash
# Edit sqldeployment.yaml to use the first example
# Then create the deployment
vvp2 deployment create -n default -f examples/sqldeployment.yaml
```

### 3. List Deployments

```bash
vvp2 deployment list -n default
```

### 4. Get Deployment Details

```bash
vvp2 deployment get sql-with-session-cluster -n default -o yaml
```

### 5. Update Deployment

```bash
# Modify the YAML file, then:
vvp2 deployment update sql-with-session-cluster -n default -f examples/test-sql-session.yaml
```

### 6. Delete Deployment

```bash
vvp2 deployment delete sql-with-session-cluster -n default
```

## Key Differences: Session Cluster vs Deployment Target

### Session Cluster
- **Pros**:
  - Faster deployment startup (cluster already running)
  - Resource sharing across multiple SQL jobs
  - Good for interactive queries and development
  - Cost-effective for multiple small jobs

- **Cons**:
  - Shared resources may impact performance
  - Jobs compete for task slots
  - Must manage session cluster lifecycle separately

- **Configuration**:
  ```yaml
  spec:
    deploymentTargetId: null
    deploymentTargetName: null
    sessionClusterName: my-sql-session
  ```

### Deployment Target
- **Pros**:
  - Dedicated resources per deployment
  - Better isolation and predictable performance
  - Independent scaling

- **Cons**:
  - Longer startup time (provisions new resources)
  - Higher resource usage
  - More expensive for many small jobs

- **Configuration**:
  ```yaml
  spec:
    deploymentTargetName: kubernetes-target
  ```

## SQL Script Guidelines

### Simple INSERT Example
```sql
INSERT INTO `mycatalog`.`db_name`.`my_table` 
VALUES ('1', 1, PROCTIME());
```

### Table Creation with Kafka Source
```sql
CREATE TABLE orders (
  order_id STRING,
  customer_id STRING,
  amount DOUBLE,
  order_time TIMESTAMP(3),
  WATERMARK FOR order_time AS order_time - INTERVAL '5' SECOND
) WITH (
  'connector' = 'kafka',
  'topic' = 'orders',
  'properties.bootstrap.servers' = 'kafka:9092',
  'format' = 'json'
);
```

### Aggregation Query
```sql
INSERT INTO order_aggregates
SELECT 
  customer_id,
  COUNT(*) as order_count,
  SUM(amount) as total_amount,
  TUMBLE_END(order_time, INTERVAL '1' HOUR) as window_end
FROM orders
GROUP BY 
  customer_id,
  TUMBLE(order_time, INTERVAL '1' HOUR);
```

## Important Configuration Fields

### Artifact Configuration
- `kind`: Must be `SQLSCRIPT` for SQL deployments
- `sqlScript`: The SQL script to execute (can be multi-line)

### State Management
- `restoreStrategy`:
  - `LATEST_STATE`: Restore from latest checkpoint/savepoint
  - `LATEST_SAVEPOINT`: Restore only from savepoints
  - `NONE`: Start fresh without state

- `upgradeStrategy`:
  - `STATEFUL`: Preserve state during upgrades
  - `STATELESS`: Restart without state

### Resource Configuration
For deployment target deployments, you can specify resources:
```yaml
resources:
  jobmanager:
    cpu: 1.0
    memory: 1024m
  taskmanager:
    cpu: 2.0
    memory: 2048m
```

For session cluster deployments, resources are inherited from the session cluster.

### Flink Configuration
Common settings:
```yaml
flinkConfiguration:
  taskmanager.numberOfTaskSlots: "2"
  execution.checkpointing.interval: "60s"
  execution.checkpointing.mode: EXACTLY_ONCE
  state.backend.type: rocksdb
  state.checkpoints.dir: "s3://bucket/checkpoints"
  state.savepoints.dir: "s3://bucket/savepoints"
```

## Troubleshooting

### Deployment Fails to Start
1. Check session cluster status:
   ```bash
   vvp2 sessioncluster get my-sql-session -n default
   ```

2. Verify deployment target exists:
   ```bash
   vvp2 deployment-target get kubernetes-target -n default
   ```

3. Check deployment logs through VVP UI or API

### SQL Syntax Errors
- Validate SQL in Flink SQL client first
- Check catalog and table names are correct
- Ensure all connectors are available in the Flink image

### Resource Issues
- For session clusters: Check available task slots
- For deployment targets: Verify Kubernetes resources
- Adjust parallelism if needed

## Testing Workflow

1. **Create session cluster** (if using session cluster approach):
   ```bash
   vvp2 sessioncluster create -n default -f examples/sessioncluster.yaml
   ```

2. **Create test deployment**:
   ```bash
   vvp2 deployment create -n default -f examples/test-sql-session.yaml
   ```

3. **Monitor deployment**:
   ```bash
   vvp2 deployment get test-sql-with-session -n default
   ```

4. **Update if needed**:
   ```bash
   # Edit the YAML file
   vvp2 deployment update test-sql-with-session -n default -f examples/test-sql-session.yaml
   ```

5. **Clean up**:
   ```bash
   vvp2 deployment delete test-sql-with-session -n default
   vvp2 sessioncluster delete my-sql-session -n default
   ```

## Additional Resources

- [Ververica Platform Documentation](https://docs.ververica.com/)
- [Apache Flink SQL Documentation](https://nightlies.apache.org/flink/flink-docs-stable/docs/dev/table/sql/overview/)
- [VVP2 CLI README](../README.md)
