# CI/CD Integration Examples

## GitHub Actions

### Basic Validation Workflow

```yaml
# .github/workflows/backstage-validate.yml
name: Validate Backstage Catalog

on:
  push:
    paths:
      - 'catalog-info.yaml'
  pull_request:
    paths:
      - 'catalog-info.yaml'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download backstage-gen-cli
        run: |
          curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
          chmod +x backstage-gen-cli

      - name: Validate catalog
        run: ./backstage-gen-cli lint --strict

      - name: Check annotations
        run: |
          ./backstage-gen-cli inspect --json | jq '.gitRemote'
```

### Generate and Commit Workflow

```yaml
# .github/workflows/generate-catalog.yml
name: Generate Backstage Catalog

on:
  workflow_dispatch:
  push:
    branches: [main]
    paths-ignore:
      - 'catalog-info.yaml'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Download backstage-gen-cli
        run: |
          curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
          chmod +x backstage-gen-cli

      - name: Generate catalog
        run: ./backstage-gen-cli generate --force

      - name: Validate generated catalog
        run: ./backstage-gen-cli lint --strict

      - name: Commit changes
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add catalog-info.yaml
          git diff --staged --quiet || git commit -m "chore: update catalog-info.yaml"
          git push
```

### Multi-Service Monorepo Workflow

```yaml
# .github/workflows/validate-monorepo.yml
name: Validate All Catalogs

on:
  push:
    paths:
      - '**/catalog-info.yaml'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download backstage-gen-cli
        run: |
          curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o /usr/local/bin/backstage-gen-cli
          chmod +x /usr/local/bin/backstage-gen-cli

      - name: Validate all catalogs
        run: |
          find . -name "catalog-info.yaml" -type f | while read file; do
            echo "Validating: $file"
            backstage-gen-cli lint -f "$file" --strict
          done
```

## GitLab CI

### Basic Validation

```yaml
# .gitlab-ci.yml
stages:
  - validate

validate-catalog:
  stage: validate
  image: alpine:latest
  before_script:
    - apk add --no-cache curl
    - curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o /usr/local/bin/backstage-gen-cli
    - chmod +x /usr/local/bin/backstage-gen-cli
  script:
    - backstage-gen-cli lint --strict
  rules:
    - changes:
        - catalog-info.yaml
```

### Generate on Merge Request

```yaml
# .gitlab-ci.yml
generate-catalog:
  stage: build
  image: alpine:latest
  before_script:
    - apk add --no-cache curl git
    - curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o /usr/local/bin/backstage-gen-cli
    - chmod +x /usr/local/bin/backstage-gen-cli
  script:
    - backstage-gen-cli generate --dry-run
    - backstage-gen-cli generate --force
    - backstage-gen-cli lint --strict
  artifacts:
    paths:
      - catalog-info.yaml
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
```

## Jenkins Pipeline

### Declarative Pipeline

```groovy
// Jenkinsfile
pipeline {
    agent any

    stages {
        stage('Setup') {
            steps {
                sh '''
                    curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
                    chmod +x backstage-gen-cli
                '''
            }
        }

        stage('Validate') {
            steps {
                sh './backstage-gen-cli lint --strict'
            }
        }

        stage('Generate') {
            when {
                branch 'main'
            }
            steps {
                sh './backstage-gen-cli generate --force'
                sh './backstage-gen-cli lint --strict'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'catalog-info.yaml', allowEmptyArchive: true
        }
    }
}
```

## CircleCI

```yaml
# .circleci/config.yml
version: 2.1

jobs:
  validate:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - run:
          name: Install backstage-gen-cli
          command: |
            curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
            chmod +x backstage-gen-cli
      - run:
          name: Validate catalog
          command: ./backstage-gen-cli lint --strict

workflows:
  version: 2
  validate:
    jobs:
      - validate:
          filters:
            branches:
              only: /.*/
```

## Azure DevOps

```yaml
# azure-pipelines.yml
trigger:
  paths:
    include:
      - catalog-info.yaml

pool:
  vmImage: 'ubuntu-latest'

steps:
  - script: |
      curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
      chmod +x backstage-gen-cli
    displayName: 'Install backstage-gen-cli'

  - script: ./backstage-gen-cli lint --strict
    displayName: 'Validate catalog'

  - script: ./backstage-gen-cli inspect --json
    displayName: 'Show inspection results'
```

## Pre-commit Hook

### Using backstage-gen-cli hooks

```bash
backstage-gen-cli hooks install
```

### Manual pre-commit config

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: backstage-lint
        name: Validate Backstage Catalog
        entry: backstage-gen-cli lint --strict
        language: system
        files: catalog-info.yaml
        pass_filenames: false
```

## Tips for CI/CD

1. **Cache the binary**: Download once and cache for faster builds
2. **Use `--quiet` flag**: Reduce noise in CI logs with `-q`
3. **Use `--strict` mode**: Treat warnings as errors in CI
4. **JSON output**: Use `--json` for programmatic processing
5. **Exit codes**: Non-zero exit on validation failure
