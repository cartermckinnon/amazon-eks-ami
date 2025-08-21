# GitHub Bot

This bot responds to comments on GitHub issues and pull requests to trigger automated workflows.

## Commands

### `/ci` - Continuous Integration (Pull Requests Only)

Triggers privileged CI jobs that require AWS credentials. This command can only be used on pull requests.

**Usage:**
```
/ci
/ci test
/ci build
```

**Named Arguments:**
You can provide additional configuration using the `+name args` syntax on subsequent lines:
```
/ci test
+workflow:timeout 60m
+os_distros al2023
```

### `/repro` - Issue Reproduction (Issues Only)

Creates a reproduction environment for customer-reported issues. This command can only be used on issues, not pull requests.

**Usage:**
```
/repro
/repro quick-test
```

**NodeConfig Options:**
You can specify NodeConfig options using the `+nodeconfig key=value` syntax:
```
/repro
+nodeconfig instance.localStorage=RAID0
+nodeconfig kubelet.maxPods=110
+nodeconfig instance.instanceType=i4i.xlarge
```

**AMI Release Tags:**
You can specify a specific AMI release tag using the `+ami` syntax:
```
/repro
+ami v20250704
+nodeconfig instance.localStorage=RAID0
```

**Example for Issue #2386:**
```
/repro
+ami v20250620
+nodeconfig instance.localStorage=RAID0
```

This would create a reproduction environment using the v20250620 AMI release with the RAID0 localStorage strategy to test the issue described in #2386.

**Advanced Example:**
```
/repro quick-test
+ami v20250704
+nodeconfig instance.localStorage=RAID0
+nodeconfig kubelet.maxPods=110
```

### `/echo` - Echo Test

Simple echo command for testing bot functionality.

**Usage:**
```
/echo hello world
```

### `/clear` - Clear Bot Comments

Removes all previous bot comments from the current issue/PR.

**Usage:**
```
/clear
```

## Authorization

Only users with `OWNER` or `MEMBER` association can use bot commands. Organization membership must be set to public for the `MEMBER` association to work.

## Implementation Notes

- The `/repro` command integrates with existing kubetest2 tooling for full reproduction workflows
- When no AMI release tag is specified, the latest recommended AMI is retrieved from AWS SSM Parameter Store
- When an AMI release tag is specified (e.g., `+ami v20250620`), the system resolves it to the actual AMI ID using EC2 DescribeImages
- NodeConfig options are passed through to kubetest2 for cluster configuration
- The bot distinguishes between issues and pull requests to enforce command restrictions
- Full integration with existing CI infrastructure (AWS credentials, log buckets, etc.)

## AMI Resolution Logic

1. **With AMI Release Tag**: Uses EC2 DescribeImages to find AMI by name pattern
2. **Without AMI Release Tag**: Uses SSM Parameter Store to get latest recommended AMI
3. **Fallback**: Builds fresh AMI if neither method succeeds (rare case)

## Future Enhancements

- Add issue-specific test focusing based on issue content analysis
- Add validation for NodeConfig option values
- Add support for custom Kubernetes versions and OS distributions
- Enhanced error reporting and debugging capabilities
