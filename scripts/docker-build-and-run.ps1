$TAG = Get-Date -Format "yyyy-MMdd-HHMMSS"
$TAG = "latest"
$IMAGE_TAG = "ubuntu:$($TAG)"


# 定义构建参数（根据需要取消注释）
$buildArgs = @(
    # "--no-cache",  # 如果不需要缓存，请取消注释
    # "--progress=plain",  # 如果需要简单的进度输出，请取消注释
    # "--platform", "linux/arm64",  # 如果需要指定平台，请取消注释
    ""
)
docker build . -t $IMAGE_TAG -f dockerfiles\Dockerfile $buildArgs


$PROJECT_NAME = $(Get-Item -Path ".").Name
# $USER_NAME = "ubuntu"
# $WORKDIR = "/home/$USER_NAME/$PROJECT_NAME"
$USER_NAME = "root"
$WORKDIR = "/$USER_NAME/$PROJECT_NAME"
# 定义运行参数（根据需要取消注释）
$buildArgs = @(
    "-i",    # -i --interactive: Keep STDIN open even if not attached
    "-t",    # -t --tty: Allocate a pseudo-TTY
    "--rm"  # --rm: Automatically remove the container and its associated anonymous volumes when it exits
    "-v", "${PWD}:$WORKDIR", # 挂载当前目录
    "-w", "$WORKDIR",      # 指定工作目录
    "-u", "$USER_NAME",    # 指定用户
    # [待验证] GPU 支持
    "--gpus", "all",       # 使用所有GPU [win] 参考 https://docs.docker.com/desktop/gpu/ 配置
    ""
)
docker run $buildArgs $IMAGE_TAG /bin/bash --login