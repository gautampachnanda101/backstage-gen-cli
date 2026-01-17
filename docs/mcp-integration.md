# MCP Integration

## Overview

MCP (Model Context Protocol) enables AI models to interact with external tools and data sources. This document covers how backstage-gen-cli can integrate with MCP-enabled systems.

## Current Status

MCP integration is planned for future releases. This document outlines the research and planned implementation.

## Related Projects

### backstage-mcp

[backstage-mcp](https://github.com/p7ayfu77/backstage-mcp) provides an MCP server that exposes Backstage catalog data to AI models.

Features:
- Query Backstage catalog entities
- Search for components, APIs, systems
- Retrieve entity details and relationships

### @mexl/backstage-plugin-catalog-backend-module-mcp

[NPM Package](https://www.npmjs.com/package/@mexl/backstage-plugin-catalog-backend-module-mcp) provides a Backstage backend module for MCP integration.

Features:
- MCP server integration in Backstage backend
- Entity registration via MCP
- Catalog updates through MCP protocol

## Planned Features

### 1. MCP Client Support

Add MCP client capability to backstage-gen-cli for pushing catalog entities directly to Backstage instances.

```bash
# Future command
backstage-gen-cli push --mcp-endpoint http://backstage:7007/mcp
```

### 2. `backstage-gen-cli push` Command

New command for pushing generated catalogs to Backstage:

```bash
# Push to Backstage via MCP
backstage-gen-cli push

# Push with custom endpoint
backstage-gen-cli push --endpoint http://backstage.internal:7007

# Dry-run to see what would be pushed
backstage-gen-cli push --dry-run

# Push and wait for registration
backstage-gen-cli push --wait
```

### 3. Configuration

```yaml
# .backstage-gen.yaml
mcp:
  enabled: true
  endpoint: http://backstage:7007/mcp
  auth:
    type: bearer
    token: ${BACKSTAGE_TOKEN}
```

## Implementation Plan

1. **Research Phase** (Current)
   - Evaluate backstage-mcp compatibility
   - Review MCP protocol specification
   - Design CLI integration approach

2. **Development Phase**
   - Add MCP client library
   - Implement `push` command
   - Add configuration options

3. **Testing Phase**
   - Integration tests with backstage-mcp
   - End-to-end testing with real Backstage instance

## Using with AI Tools

### Claude Desktop Integration

backstage-gen-cli can be used alongside MCP servers in Claude Desktop:

```json
{
  "mcpServers": {
    "backstage": {
      "command": "npx",
      "args": ["-y", "backstage-mcp"],
      "env": {
        "BACKSTAGE_URL": "http://localhost:7007"
      }
    }
  }
}
```

Then use backstage-gen-cli to generate catalogs that can be queried via the MCP server.

### Workflow Example

1. Generate catalog with backstage-gen-cli
2. Commit and push to repository
3. Backstage discovers and registers entity
4. Query entity via MCP-enabled AI tools

## Resources

- [MCP Specification](https://modelcontextprotocol.io/)
- [backstage-mcp GitHub](https://github.com/p7ayfu77/backstage-mcp)
- [Backstage Software Catalog](https://backstage.io/docs/features/software-catalog/)

## Contributing

If you'd like to help implement MCP integration:

1. Check [GitHub Issues](https://github.com/gautampachnanda101/backstage-gen-cli/issues) for MCP-related tasks
2. Review the planned features above
3. Submit PRs with implementation or suggestions
