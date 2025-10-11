# 使用BUILDPLATFORM和TARGETPLATFORM进行多架构构建
# FROM --platform=$BUILDPLATFORM golang:1.24 AS golang-builder

# 设置目标平台环境变量
ARG TARGETOS
ARG TARGETARCH

ARG GO_VERSION=1.23
ARG UV_VERSION=0.8.8
ARG NODE_IMAGETAG=24

FROM golang:${GO_VERSION} AS golang-builder
FROM ghcr.io/astral-sh/uv:${UV_VERSION} AS uv-builder
FROM node:${NODE_IMAGETAG} AS node-builder

FROM ubuntu:24.04

# ============================================================
# 「必须」软件列表 MUST_INSTALL_LIST
# ============================================================
ARG MUST_INSTALL_LIST="ubuntu-server"
ARG MUST_INSTALL_LIST="$MUST_INSTALL_LIST ca-certificates locales"
ARG MUST_INSTALL_LIST="$MUST_INSTALL_LIST zsh"

# ============================================================
# 「可选」软件列表 OPTION_INSTALL_LIST
# ============================================================
ARG OPTION_INSTALL_LIST=""
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST vim nano tree pciutils iputils-ping dnsutils"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST git openssh-client"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST build-essential"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST cmake ninja-build clangd"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST python3-dev"
# ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST golang"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST postgresql-client"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST redis-tools"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST mysql-client"
ARG OPTION_INSTALL_LIST="$OPTION_INSTALL_LIST lua5.3"

# 自定义镜像源
# ARG MIRRORS_URL="mirrors.ustc.edu.cn"
ARG MIRRORS_URL="mirrors.tencentyun.com"

# 是否清理 apt 缓存
ARG CLEAN_APT_CACHE=0

# 安装软件
ENV DEBIAN_FRONTEND=noninteractive
ENV INSTALL_LIST="$MUST_INSTALL_LIST $OPTION_INSTALL_LIST"
# set -x 选项会使得每个命令在执行前都会被打印出来，报错时会显示是哪个命令出错
RUN set -x && \
    sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' ${HOME}/.bashrc && \
    sed -i -e "s#//.*archive.ubuntu.com#//${MIRRORS_URL}#g" -e "s#//ports.ubuntu.com#//${MIRRORS_URL}#g" /etc/apt/sources.list.d/ubuntu.sources && \
    apt-get update && apt-get upgrade -y && \
    apt-get install -y --no-install-recommends ${INSTALL_LIST} && \
    apt-get autoremove -y && \
    if [ $CLEAN_APT_CACHE -ne 1 ]; then { apt-get clean all; rm -rf /var/lib/apt/lists/*; } fi && \
    localedef -i zh_CN -c -f UTF-8 -A /usr/share/locale/locale.alias zh_CN.UTF-8 && \
    rm -rf /tmp/*

ENV LANG=zh_CN.utf8

# 系统时间错误可能导致: `tls: failed to verify certificate: x509: certificate signed by unknown authority`
# 修改时区 (date -R 检查当前时间)
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 启用 bash 补全
# RUN sed -i.bak '/^#if ! shopt -oq posix; then$/,/^#fi$/s/^#//' /etc/bash.bashrc

# 设置 zsh 为默认 shell
RUN set -x && \
    chsh -s /usr/bin/zsh && \
    # ====== install oh-my-zsh (official) ======
    # git clone https://github.com/ohmyzsh/ohmyzsh.git --depth 1 /tmp/ohmyzsh && cd /tmp/ohmyzsh/tools && \
    # sh install.sh && rm -rf /tmp/ohmyzsh && \
    # git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-/root/.oh-my-zsh/custom}/plugins/zsh-autosuggestions && \
    # ====== install oh-my-zsh (mirror) ======
    git clone http://mirrors.tuna.tsinghua.edu.cn/git/ohmyzsh.git --depth 1 && cd ohmyzsh/tools && \
    REMOTE=https://mirrors.tuna.tsinghua.edu.cn/git/ohmyzsh.git sh install.sh && \
    git clone https://gitee.com/henryzhuhr/zsh-autosuggestions.git --depth 1 ${ZSH_CUSTOM:-/root/.oh-my-zsh/custom}/plugins/zsh-autosuggestions && \
    # ==============================
    # set zsh theme and plugins
    echo "ZSH_THEME='candy'" >> ~/.zshrc && \
    echo "plugins=(macos dotenv git cp python bundler rake rbenv ruby zsh-autosuggestions)" >> ~/.zshrc && \
    echo "ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE='fg=#00BFFF,bg=gray,bold,underline'" >> ~/.zshrc && \
    echo "ZSH_AUTOSUGGEST_STRATEGY=(history completion)" >> ~/.zshrc && \
    echo "source \$ZSH/oh-my-zsh.sh" >> ~/.zshrc && \
    echo "alias ll='ls -alh'" >> ~/.zshrc && \
    echo "alias la='ls -A'" >> ~/.zshrc && \
    echo "alias l='ls -CF'" >> ~/.zshrc && \
    echo "alias h='history'" >> ~/.zshrc && \
    echo "alias mkdir='mkdir -p'" >> ~/.zshrc && \
    echo "alias grep='grep --color=auto'" >> ~/.zshrc && \
    echo "alias egrep='egrep --color=auto'" >> ~/.zshrc && \
    echo "alias fgrep='fgrep --color=auto'" >> ~/.zshrc && \
    rm -rf /tmp/*

ARG USER_NAME="root"
ARG USER_HOME="/${USER_NAME}"

# ============================================================
#   ⬇️ 扩展 Dockerfile
# ============================================================


# ============================================================
#   Golang 环境 - 支持多架构
# ============================================================
ARG TARGETOS
ARG TARGETARCH
COPY --from=golang-builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"
ENV PATH="$USER_HOME/go/bin:$PATH"
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

# ============================================================
#   Python 环境
# ============================================================
# # PEP 668 – Python base environments Python 增强提案 (PEP)：https://realpython.com/python-virtual-environments-a-primer/?ref=yaolong.net
# # https://www.adamsdesk.com/posts/resolve-pip-externally-managed-environment
# # 临时解决方案  --ignoreExternallyManagedEnvironment
# 使用 uv
COPY --from=uv-builder /uv /uvx /bin/

# 默认使用系统 Python 环境
ENV UV_SYSTEM=true 
ENV VIRTUAL_ENV=$USER_HOME/.venv

ADD . $USER_HOME/tmp
ENV UV_DEFAULT_INDEX="https://mirrors.cloud.tencent.com/pypi/simple/"
# ENV UV_INDEX="https://mirrors.ustc.edu.cn/pypi/simple"
RUN set -x && \
    uv venv $VIRTUAL_ENV && \
    echo "source ${VIRTUAL_ENV}/bin/activate" >> $USER_HOME/.zshrc && \
    rm -rf $(uv cache dir) && \
    rm -rf $USER_HOME/tmp && \
    # 在 uv 没有支持全局变量前，临时加入 history 文件
    echo ": $(date +%s):0;uv add  --active" >> $USER_HOME/.zsh_history && \
    echo ": $(date +%s):0;uv run  --active" >> $USER_HOME/.zsh_history && \
    echo ": $(date +%s):0;uv sync --active" >> $USER_HOME/.zsh_history && \
    echo "finished" && \
    rm -rf /tmp/*

# 环境初始化脚本，加快构建镜像的速度
RUN echo ": $(date +%s):0;zsh scripts/init.sh" >> $USER_HOME/.zsh_history

# 提升虚拟环境 python 的优先级
ENV PATH=${VIRTUAL_ENV}/bin:$PATH
ENV VIRTUAL_ENV_DISABLE_PROMPT=1

ADD scripts/init.sh /root/scripts/init.sh

