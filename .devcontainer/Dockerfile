ARG GO_VERSION=1.23
ARG GO_IMAGE=mcr.microsoft.com/devcontainers/go
ARG ALPINE_VERSION=bookworm

FROM ${GO_IMAGE}:${GO_VERSION}-${ALPINE_VERSION}

# Install dependencies
RUN apt update \
    && apt install -yqq jq \
    && apt clean -yqq

# Update non-root user
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Switch to non-root user
USER $USERNAME
WORKDIR /workspace

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest