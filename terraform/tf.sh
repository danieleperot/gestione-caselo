#!/usr/bin/env bash

function help() {
    echo "Usage: $0 {state} {environment} [terraform commands...]" 1>&2
    echo -e "\tNote: for \"global\" state, \"environment\" must be skipped." 1>&2
    exit 1
}

terraform_dir="$(git rev-parse --show-toplevel)/terraform"
env_path="${terraform_dir}/.env"

if [ ! -f $env_path ]; then
    echo "Missing .env file in ${terraform_dir}. Create it by copying \`.env.example\` and filling it with data appropriate to your setup." 1>&2
    exit 1
fi

set -o allexport && source $env_path && set +o allexport


if [ -z $1 ]; then
    echo "ERROR: no \"state\" parameter provided." 1>&2
    help
fi

state=$1
state_dir="$terraform_dir/states/$state"
if [ ! -d $state_dir ]; then
    echo "ERROR: invalid state \"$state\" provided - the directory $state_dir does not exist!"
    help
fi

shift 1

tf_region="eu-south-1" # TODO: get from env file / state?

var_flags=""
var_flags="$var_flags -var=region=$tf_region"
var_flags="$var_flags -var=project=$PROJECT_NAME"

tf_env=""
tf_vars_file=""
if [ "$state" != "global" ]; then
    if [ -z $1 ]; then
        echo "ERROR: no \"environment\" parameter provided. It is required for every state (except \"global\")." 1>&2
        help
    fi

    tf_env="$1"
    tf_vars_file="$terraform_dir/tfvars/$tf_env.tfvars"
    if [ ! -f $tf_vars_file ]; then
        echo "ERROR: invalid environment \"$tf_env\" provided - the file $tf_vars_file does not exist!"
        help
    fi

    var_flags="$var_flags -var=environment=$tf_env"
    shift 1
fi


command=$1
shift 1
if [ "$command" == "init" ]; then
    flags="$flags -backend-config=bucket=$STATE_BUCKET"
    flags="$flags -backend-config=region=$tf_region"
    if [ "$state" != "global" ]; then
        flags="$flags -backend-config=key=$tf_env/$state.tfstate"
    fi
else
    if [ "$state" != "global" ]; then
        var_flags="$var_flags -var-file=$tf_vars_file"
    fi
fi

if [ "$command" != "state" ]; then
    flags="$flags $var_flags"
fi

set -x
TF_DATA_DIR="$state_dir/.${tf_env:-global}.terraform" terraform -chdir=$state_dir $command $flags $@
