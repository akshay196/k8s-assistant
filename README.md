# k8s-assistant: Your AI-Powered Kubernetes Assistant

## Setup

Obtain an API key from https://platform.openai.com/api-keys. Provide
API key using environment variable `OPENAI_API_KEY`.

Start k8s-assistant:

```bash
go run main.go
```

## Roadmap

- [ ] Accept flags to customize app behaviour.
- [ ] (UX) Stream executing command output to stdout instead of printing all at once.
