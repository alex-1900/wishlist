FROM mcr.microsoft.com/devcontainers/go:1.25
USER root

RUN apt update && apt install -y \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
  && apt install -y nodejs \
  && rm -rf /var/lib/apt/lists/* \
  && chsh -s $(which zsh) vscode

RUN mkdir -p /usr/local/share/npm-global && \
  chown -R vscode:vscode /usr/local/share


RUN mkdir -p /home/vscode/.claude && \
  chown -R vscode:vscode /home/vscode/.claude

USER vscode

ENV DEVCONTAINER=true
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

ENV SHELL=/bin/zsh
ENV EDITOR=nano
ENV VISUAL=nano

RUN npm install -g @anthropic-ai/claude-code@latest
