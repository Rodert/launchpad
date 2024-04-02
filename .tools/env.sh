#!/usr/bin/env bash

set -eux

set_basic_env() {
    export PATH=$PATH:$(go env GOPATH)/bin

    export GIT_BRANCH_BASENAME=$(git rev-parse --abbrev-ref HEAD)
    export GIT_COMMIT=$(git rev-parse HEAD)
    export GIT_COMMIT_SHORT=$(git rev-parse HEAD)
    export GIT_MESSAGE=$(git log --pretty=format:"%s" | head -n 1)
    export GIT_COMMITTER=$(git log --pretty=format:"%cn" | head -n 1)

    if [[ "${GIT_BRANCH_BASENAME}" = "main" ]];then
        export PROJECT_ENV="release"
    elif [[ "${GIT_BRANCH_BASENAME}" = "preview" ]];then
        export PROJECT_ENV="debug"
    else
        export PROJECT_ENV="test"
    fi

    if [[ "${GIT_BRANCH_BASENAME}" =~ ^fix|feature\/.*-([0-9]{7,\}) ]]; then
        export USER_FEISHU_ID=${BASH_REMATCH[1]}
    fi

    export PROJECT_NAME=`basename $(git rev-parse --show-toplevel)`
    export SCRIPT_DIR=$(dirname $0)
    export PROJECT_ROOT="${SCRIPT_DIR}/../"
    export BUILD_DIR="${PROJECT_ROOT}/.build"
    export BUILT_TIME="$(date "+%Y-%m-%d %H:%M:%S")"
    export VERSION="$(cat version)"

    # 阿里云云效通过将环境变量输出到 .env 文件中后任务内才可直接使用该环境变量。
    # 注意: 阿里云云效 .env 文件中的环境变量需要以 USER_ 开头。
    export USER_GIT_COMMITTER="${GIT_COMMITTER}"
    export USER_GIT_BRANCH_BASENAME="${GIT_BRANCH_BASENAME}"
    export GOPROXY=https://goproxy.cn,direct
    echo "$(env)" > .env
}

set_basic_env