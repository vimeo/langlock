#!/bin/bash

if [ ! "${BASH_SOURCE[0]}" -ef "$0" ]
then
    echo "ERROR: Hey, you should execute this script directly instead of sourcing it!"
    return 1
fi

RELATIVE_SCRIPT_PATH="${BASH_SOURCE%/*}"
pushd "$RELATIVE_SCRIPT_PATH" >/dev/null
SCRIPT_PATH="$(pwd -P)" # Absolute script path
popd >/dev/null
readonly SCRIPT_PATH

langlock_tarballname="langlock-pre-receive-env.tar.gz"
langlock_imagename="langlock-env"
langlock_containername="langlock-env"
langlock_build_dir="build"
langlock_local_scanner_bin_path="$SCRIPT_PATH/$langlock_build_dir/detectcreds.bin"
readonly langlock_tarballname
readonly langlock_imagename
readonly langlock_containername

RED='\033[31m'
GREEN='\033[92m'
BLUE='\033[96m'
BOLD='\033[1m'
PLAIN='\033[0m'
readonly RED
readonly GREEN
readonly BLUE
readonly BOLD
readonly PLAIN

USAGE=$(cat <<-END
Usage:
    ./langlock.sh SUBCOMMAND

Subcommands:
    help, -h, --help
        Print this usage message.

    build SLACK_WEBHOOK_URL
        Build the pre-receive hook environment.

    build-local
        Build the GoLang secret scanner for local use

    deploy IP_ADDRESS_OF_GITHUB_ENTERPRISE_SERVER
        Deploy the pre-receive hook environment to the server.

    hash STRING_TO_HASH
        Generate the hash for a given string (as would be used for the key in the allow list file).

    test-scanner
        Invoke the tests of the golang module locally.

    test-e2e PATH_TO_TEST_TARGET_REPO
        Invoke the end-to-end tests on any target repo.

    scan ARGS_FOR_LOCAL_SCANNER...
        Run the scanner locally

    scan-dir PATH_TO_TARGET_DIR ARGS_FOR_LOCAL_SCANNER...
        Run the scanner locally against a directory recursively. Skips non-text files.

    scan-commits PATH_TO_REPO COMMIT_HASHES ARGS_FOR_SCANNER...
        Run the scanner locally against specific commits in a specific repo. The commits
        must be space-separated.
END
)
readonly USAGE

function error() {
    echo -e "${PLAIN}${RED}${BOLD}ERROR: $1${PLAIN}"
    exit 1
}


# Ensure that Langlock's go code is on the gopath
if [[ $GOPATH != *"$SCRIPT_PATH"* ]];then
    export GOPATH="$SCRIPT_PATH:$GOPATH"
fi


function print_usage() {
  echo -e "$USAGE"
}

function get_date () {
  date "+%Y%m%d_%H%M%S"
}

function build_localscanner () {
    echo "Building Langlock's secret scanner locally"
    # Only build the binary if it does not already exist, or if the
    # "replace" flag (-r) is used
    if [[ ( ! -f "$langlock_local_scanner_bin_path" || "$1" = "-r" ) ]]; then
      (cd "${SCRIPT_PATH}/src/detectcreds" && go build -o "$langlock_local_scanner_bin_path")
    else
      echo "Building Langlock's secret scanner skipped since binary already exiss"
    fi
}

function hash() {
    openssl sha256 <(echo "$1") | cut -d ' ' -f 2
}

function build_hookenv() {
    echo "Building Langlock tarball for hook environment."
    local slack_webhook_url="$1"
    if [ -z "$slack_webhook_url" ]; then
        error "Must provide slack webhook url as first argument to deploy command"
    fi
    docker build -f "${SCRIPT_PATH}/Dockerfile" --build-arg SLACK_WEBHOOK_URL="$slack_webhook_url" -t "$langlock_imagename" "$SCRIPT_PATH" || error "Problem running docker build"
    docker container rm "$langlock_containername" 2>/dev/null || true;
    docker create --name "$langlock_containername" "$langlock_imagename" /bin/true || error "Problem creating a new container"
    mkdir -p "${SCRIPT_PATH}/$langlock_build_dir"
    docker export "$langlock_containername" | gzip > "${SCRIPT_PATH}/${langlock_build_dir}/${langlock_tarballname}" || error "Problem creating the tarball"
    # docker container rm "$langlock_containername" 2>/dev/null || true;
    echo -e "\n\n${GREEN}${BOLD}SUCCESS:${PLAIN} created the langlock pre-receive hook environment tarball (${SCRIPT_PATH}/${langlock_build_dir}/${langlock_tarballname})"
    echo -e "\n${BLUE}In order to use Langlock server-side, you must next deploy the tarball that you just created to your GitHub Enterprise Server as a pre-receive hook environment. You may use the ${PLAIN}${BOLD}./langlock deploy${PLAIN}${BLUE} command to do so. Please see the README for further details.${PLAIN}"
}

function deploy_hookenv() {
    echo "Deploying Langlock hook environment to github enterprise: $1"
    echo "This may take a few minutes."
    local github_enterprise_ip_address="$1"

    if [ -z "$github_enterprise_ip_address" ]; then
        error "Must provide github enterprise ip address as second argument to deploy command"
    fi

    scp -P 122 "${SCRIPT_PATH}/${langlock_build_dir}/${langlock_tarballname}" "admin@${github_enterprise_ip_address}:/home/admin" || error "Problem uploading the tarball to the GitHub Enterprise Server ip address"

    deployed_env_name="langlock-env-$(get_date)"
    ssh -p 122 "admin@$github_enterprise_ip_address" "ghe-hook-env-create $deployed_env_name /home/admin/${langlock_tarballname}" || error "Problem SSHing into the GitHub Enterprise Server and running the ghe-hook-env-create command"
    echo -e "\n\n${GREEN}${BOLD}SUCCESS:${PLAIN} Created new pre-receive hook environment named ${deployed_env_name} in GitHub Enterprise Server."
    echo -e "\n${BLUE}Note that, as a final step, you must next use your browser to configure the pre-receive hook to use the version of the environment that you just uploaded. Please consult the README for more details.${PLAIN}"
    return 0
}

function test_scanner() {
    echo "Running unit tests for detectcreds golang tool"
    (cd "${SCRIPT_PATH}/src/detectcreds" && go test ./...)
}

function test_e2e() {
    echo "Running end-to-end tests for Langlock on repo: $1"
    local path_to_test_repo="$1"
    if [ -z "$path_to_test_repo" ]; then
        error "Must provide the path of the repo (with Langlock enabled) on which to perform the end-to-end tests"
    fi
    (cd "$path_to_test_repo" && bash "$SCRIPT_PATH/test-e2e.sh")
}

function scan() {
  #pushd "$SCRIPT_PATH/src/detectcreds/" >/dev/null
  build_localscanner >/dev/null
  $langlock_local_scanner_bin_path $@

  #go run "$SCRIPT_PATH/src/detectcreds/main.go" $@
  #popd >/dev/null
}

function scanDir() {
  build_localscanner >/dev/null
  target_dir="$1"
  shift
  find "$target_dir" -type f -not -path "*/.git/*" -exec grep -I -q . {} \; -print | xargs -L1 -I % sh -c "$langlock_local_scanner_bin_path % $@"
}

function scan_commits() {
  path_to_target_repo="$1"
  commits_to_scan="$2"
  extra_args="$3"
  (cd $path_to_target_repo && git show "$commits_to_scan" && go run "$SCRIPT_PATH/src/detectcreds/main.go" <(git show $commits_to_scan -m --no-prefix --ignore-space-change --unified=0 -- . ":(exclude)$ALLOW_LIST_PATH" | grep -v '^@@' | egrep -v '^-[^-]') $extra_args)
}


while getopts ":h" opt; do
  case ${opt} in
    h )
      print_usage
      exit 0
      ;;
   \? )
     echo "Invalid Option: -$OPTARG" 1>&2
     exit 1
     ;;
  esac
done
shift $((OPTIND -1))


subcommand=$1; shift

if [[ -z "$subcommand" ]]; then
    print_usage
    exit 0
fi

case "$subcommand" in
  help)
    print_usage
    exit 0
    ;;
  --help)
    print_usage
    exit 0
    ;;
  build)
    arg_slack_webhook_url="$1"; shift
    build_hookenv "$arg_slack_webhook_url"
    exit 0
    ;;
  hash)
    hash "$1"
    exit 0
    ;;
  build-local)
    build_localscanner -r
    exit 0
    ;;
  deploy)
    arg_github_enterprise_ip_addr="$1"; shift
    deploy_hookenv "$arg_github_enterprise_ip_addr"
    exit 0
    ;;
  test-scanner)
    test_scanner
    exit 0
    ;;
  test-e2e)
    arg_path_to_dummy_repo_for_e2e_tests="$1"; shift
    test_e2e "$arg_path_to_dummy_repo_for_e2e_tests"
    exit 0
    ;;
  test)
    echo "Please run either:"
    echo "   • ./langlock.sh test-e2e ——— end-to-end test on a remote deployed version of Langlock in a dummy repo of your choosing or"
    echo "   • ./langloch.sh test-scannner ——— locally test Langlock's secret-detection logic."
    exit 0
    ;;
  scan)
    scan "$@"
    ;;
  scan-dir)
    scanDir "$@"
    ;;
  scan-commits)
    path_to_target_repo="$1"; shift
    commits_to_scan="$1"; shift
    extra_args="$@"
    scan_commits "$path_to_target_repo" "$commits_to_scan" "$extra_args"
    ;;
  * )
    echo "Invalid command provided"
    print_usage
    exit 1
    ;;
esac
