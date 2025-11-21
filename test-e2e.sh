#!/bin/bash


# ----------------------------------------------------------------------------
# Required manual setup
# ---------------------------------------------------------------------------

# The test cases related the enrolled users require that we there be at least
# two users with write access to the org/repos within Github Enterprise.

# 1. Please set up two Github Enterprise account
# 2. In the web UI for GitHub Enterprise, add a separate SSH key for each user
# 3. Update your ~/.ssh/config file to contain the following (indentation matters):
#
#        Host langlock_demo_main
#          HostName IP_OF_GITHUB_ENTERPRISER_SERVER
#          IdentityFile PATH_TO_PRIVATE_KEY_FILE_OF_FIRST_USER
#          User git
#          IdentitiesOnly yes
#
#        Host langlock_demo_alt
#          HostName IP_OF_GITHUB_ENTERPRISER_SERVER
#          IdentityFile PATH_TO_PRIVATE_KEY_FILE_OF_SECOND_USER
#          User git
#          IdentitiesOnly yes
#
# 4. For the target test repo, update the .git/config file to include the following:
#
#        [remote "origin"]
#                url = git@langlock_demo_main:sampleorg/demo-repo-1.git
#                fetch = +refs/heads/*:refs/remotes/origin/*
#        [remote "altorigin"]
#                url = git@langlock_demo_alt:sampleorg/demo-repo-1.git
#                fetch = +refs/heads/*:refs/remotes/origin/*
#
# 5. Start an ssh-agent session. Just type: `ssh-agent` then copy the output into
#    the terminal. If a session had already been started, use `ssh-add -D` to remove
#    any previously added keys.
# 6. Specify the specific/different origin (e.g., `git push altorigin`)
#    in order to effectively push as a different user

# ----------------------------------------------------------------------------
# Util functions
# ---------------------------------------------------------------------------

get_date () {
  date "+%Y%m%d_%H%M%S"
}

function fail() {
    echo -e "FAILED: $1"
    exit 1
}

function pass() {
    echo -e "PASSED: $1"
}


# ----------------------------------------------------------------------------
# Prepare testing environment
# ---------------------------------------------------------------------------

BRANCH_NAME="langlock_test_$(get_date)"
git checkout -b "$BRANCH_NAME" 2>/dev/null
if [ $? -ne 0 ]; then
    fail "Error when setting up testing environment -- Could not checkout new branch (${BRANCH_NAME})."
fi
git push -u origin "$BRANCH_NAME" --push-option=SKIP_LANGLOCK >/dev/null 2>&1
if [ $? -ne 0 ]; then
    fail "Error when setting up testing environment -- Could not push new branch (${BRANCH_NAME}) upstream."
fi

BASE_FILENAME="langlock-test"

NORMAL_FILE="${BASE_FILENAME}-foo.txt"
DANGEROUS_NAME_FILE="id_rsa"
BIG_FILE="${BASE_FILENAME}-big.txt"
ALLOW_LIST_FILE="langlock.json"
LOCAL_BIG_FILE_2_4MB="/usr/share/dict/web2"
LOCAL_BIG_FILE_0_5MB="${BASE_FILENAME}-local-big.txt"
if [ -f "$LOCAL_BIG_FILE_0_5MB" ]; then
  rm LOCAL_BIG_FILE_0_5MB
fi
touch LOCAL_BIG_FILE_0_5MB
head -n 48121 $LOCAL_BIG_FILE_2_4MB > $LOCAL_BIG_FILE_0_5MB
TEMP_FILE="${BASE_FILENAME}-temp.txt"

if [ -f "$NORMAL_FILE" ]; then
    git rm -f "$NORMAL_FILE" >/dev/null && git commit -m 'Remove normal file' >/dev/null && git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi
if [ -f "$DANGEROUS_NAME_FILE" ]; then
    git rm -f "$DANGEROUS_NAME_FILE" >/dev/null && git commit -m 'Remove file with suspicious/dangerous name' >/dev/null && git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi
if [ -f "$BIG_FILE" ]; then
    git rm -f "$BIG_FILE" >/dev/null && git commit -m 'Remove big file' >/dev/null && git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi

function removeAllowList {
    if [ -f "$ALLOW_LIST_FILE" ]; then
        git rm -f "$ALLOW_LIST_FILE" >/dev/null && git commit -m 'Remove allow list file' >/dev/null && git push --push-option=SKIP_LANGLOCK 2>/dev/null
    fi
}

removeAllowList

FAKE_AWS_ACCESS_KEY='aws_access_key_id=AKIAW5JVYFDT5ZB7FAKE'
FAKE_AWS_ACCESS_KEY_SHORT='AKIAW5JVYFDT5ZB7FAKE'

FAKE_AWS_ACCESS_KEY_2='aws_access_key_id=AKIAW5JVYFDT5ZB8FAKE'
FAKE_AWS_ACCESS_KEY_2_SHORT='AKIAW5JVYFDT5ZB8FAKE'


# ----------------------------------------------------------------------------
# Test cases
# ---------------------------------------------------------------------------

TEST_NAME='should allow benign pushes that do not contain secrets'
if true; then
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text"
    fi
    pass "$TEST_NAME"
fi


TEST_NAME='should block pushes that contain one suspected secret'
if true; then
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: detected secrets in branch" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep 'If the string is indeed sensitive' >/dev/null; then
        fail "$TEST_NAME —- missing remediation information in output for blocked push. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi


TEST_NAME='should block pushes that contain many suspected secrets'
if true; then
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with Fake AWS access key' >/dev/null
    echo '$password = "AmiMqqOqqqqOqQ8t";' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with generic base62 password' >/dev/null
    echo -e 'fastly 28w5hle0koSw8q-sh888q8888X7sFAKE\nGCP API key\nAIzaSyBuc29gxhCHltdsvLMSe88888QBxymFAKE\n\nhttps://hooks.slack.com/services/T1NT18888/B018888DJAB/fwMqpIUoLjLWqw88888cFAKE\n\ncurl -u thisIsAUsername:imAPassword http://example.com\n\ncurl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \n     -H "Authorization: Bearer zzzz88ZZ_z88zzzzzZZZZ888ZZZZZ_zZZzz88zZZ" \n\     -H "Content-Type:application/json"\n\nFake Cloudflare global API KEY 1b8b888888e1f6ca0c2ec297bef8bbcf88ed8\nFake Cloudflare origin CA key: v1.0-d888888cf2bd888888dcf8e6-888888cb0b9b88f8c8888a8888c88ed88fbaab88888888f888888888b360de88e888ee0fe8888f888888a8aacc8b88888f888a9caa88a88b8cae8f8a8d88ea8888888888db8ad28888\nFake HackerOne token: L888SYHhQQQqQQQ88888qqqq1mt/8QQ8jQQQQq8QqQQ=\nFake atlassian secret = ATAAT3QQQQQQQQQQQQQ88888888888888888QQ_8888888888888qqqqqqqqqqqqqqqqqqqqqqqqqq-88888888888-qqqqqqqqqq888888888QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ_8888888888888888888888888_88888QQQQQQQ=8
8888888' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with fake atlassian api token, hackerone token, cloudflare origin ca key, cloudflare global api key, cloudflare api token, curl with credentials, gcp api key, fastly credential, etc.' >/dev/null
    echo '$this->db = "postgresql://other:fakeDbPassword@localhost/otherdb?connecttime=10&application_name=myApp"' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with postgresql database connection string' >/dev/null
    echo '$combination->ints[$type] = "AIzaSyB9V5n21zmPBaIuO1SHoxZHnipZb76Azz0";' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with generic base62 password' >/dev/null
    echo 'a' >> "$DANGEROUS_NAME_FILE" && git add "$DANGEROUS_NAME_FILE" && git commit -m 'Make commit with sensitive file name' >/dev/null
    echo 'api_key_2 = "0xce572dc620c806c230eb68fbd371d1"' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with lowercase hex credential' >/dev/null
    echo 'api_key = "0xCE572DC620C806C230EB68FBD371D1"' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with uppercase hex credential' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    numSecretsFound="$(echo "$stdout" | grep "Potential secret: " | wc -l)"
    expectedNumSecrets="17"
    if [ $numSecretsFound -ne $expectedNumSecrets ]; then
        fail "$TEST_NAME -- expected to detect $expectedNumSecrets credentials, but only found ${numSecretsFound}. Stdout:\n $stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi


TEST_NAME='should block pushes that contain dangerous file names'
if true; then
    echo "a" >> "$DANGEROUS_NAME_FILE" && git add "$DANGEROUS_NAME_FILE" && git commit -m 'Make commit with a dangerous file name' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: detected secrets in branch" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$DANGEROUS_NAME_FILE" >/dev/null; then
        fail "$TEST_NAME —- missing plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep 'If the string is indeed sensitive' >/dev/null; then
        fail "$TEST_NAME —- missing remediation information in output for blocked push. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi


TEST_NAME='should be able to push secret successfully after running the provided one-liner to update the allow list'
if true; then
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    echo '{}' > "$ALLOW_LIST_FILE"
    git add "$ALLOW_LIST_FILE"; git commit -m "Clear allow list" >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code before remediation"
    fi
    remediation_command="$(echo "$stdout" | grep "Update allow list" | sed 's/^remote: *//g' | sed 's/git push//g')"
    ( eval $remediation_command ) >/dev/null 2>/dev/null << USER_INPUT_HERE_DOC
Adding this credential to the allow list because it is a false positive
USER_INPUT_HERE_DOC
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code after remediation"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text"
    fi
    pass "$TEST_NAME"
    removeAllowList
fi

TEST_NAME='should be able to push secret successfully after adding file path to allow list'
echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
echo '{}' > "$ALLOW_LIST_FILE"
git add "$ALLOW_LIST_FILE"; git commit -m "Clear allow list" >/dev/null
stdout="$(git push 2>&1)"
if [ $? -ne 1 ]; then
    fail "$TEST_NAME —- wrong status code before remediation"
fi
printf '{ "allowPaths": [{"regex": "%s", "reason": ""}]}' "langlock-test-foo.txt" > "$ALLOW_LIST_FILE"
git add "$ALLOW_LIST_FILE"; git commit -m "Add path to allow list" >/dev/null
stdout="$(git push 2>&1)"
echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
stdout="$(git push 2>&1)"
if [ $? -ne 0 ]; then
    fail "$TEST_NAME —- wrong status code after remediation"
fi
if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
    fail "$TEST_NAME —- missing correct stdout text"
fi
pass "$TEST_NAME"


TEST_NAME='should block pushes if the allow list has an incorrect format'
if true; then
    # Since the allow list is a json file, appending a single character to the end
    # of the file will make its format invalid
    echo "blah" > "$ALLOW_LIST_FILE" && git add "$ALLOW_LIST_FILE" && git commit -m 'Make commit with misformatted allow list' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: Your push was blocked because we detected invalid syntax" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "Allow list must be JSON, with the format:" >/dev/null; then
        fail "$TEST_NAME —- missing instructions for properly formatting the JSON allow list. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    removeAllowList
fi

TEST_NAME='should skip all scanning if the bypass flag is provided in the push'
if true; then
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    stdout="$(git push --push-option=SKIP_LANGLOCK 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "Server-side secret scanner (Langlock) bypassed." >/dev/null; then
        fail "$TEST_NAME —- missing bypass message in stdout text. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
fi

TEST_NAME='should be performant. i.e., should be able to process large push (~ 0.5MB) without timing out'
# ... when GitHub Enterprise Server is hosted on Google Cloud Platform (GCP) n1-standard-4 (Intel Haswell, 4 vCPUs, 15 GB memory, local SSD, max egress 10 Gbps). Note that pre-receive hooks in GHE Server are required to complete in < 5 seconds.'
if true; then
    cat "$LOCAL_BIG_FILE_0_5MB" >> "$BIG_FILE" && git add "$BIG_FILE" && git commit -m 'Make large commit' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text -- stdout:\n $stdout"
    fi
    if echo "$stdout" | grep "Server-side secret scanner (Langlock) timed out, but git push still accepted." >/dev/null; then
        fail "$TEST_NAME —- unexpected graceful timeout message in stdout text."
    fi
    pass "$TEST_NAME"
fi

TEST_NAME='should allow a large push that times out'
if true; then
    cat "$LOCAL_BIG_FILE_2_4MB" >> "$BIG_FILE" && git add "$BIG_FILE" && git commit -m 'Make large commit' >/dev/null
    cat "$LOCAL_BIG_FILE_2_4MB" >> "$BIG_FILE" && git add "$BIG_FILE" && git commit -m 'Make large commit' >/dev/null
    cat "$LOCAL_BIG_FILE_2_4MB" >> "$BIG_FILE" && git add "$BIG_FILE" && git commit -m 'Make large commit' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "Server-side secret scanner (Langlock) timed out, but git push still accepted." >/dev/null; then
        fail "$TEST_NAME —- missing graceful timeout message in stdout text. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
fi


TEST_NAME='should work when pushing multiple commits on the same branch'
if true; then
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 1' >/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 2' >/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 3' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text"
    fi
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 3' >/dev/null
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 3' >/dev/null
    echo "$FAKE_AWS_ACCESS_KEY_2" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit 3' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: detected secrets in branch" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_2_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing second plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep 'If the string is indeed sensitive' >/dev/null; then
        fail "$TEST_NAME —- missing remediation information in output for blocked push. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi


TEST_NAME='should work if pushing no commits'
if true; then
    stdout="$(git push origin "${BRANCH_NAME}:${BRANCH_NAME}_2" 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing graceful timeout message in stdout text. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
fi


TEST_NAME='should detect secrets if multiple branches are pushed'
if true; then
    stdout="$(git push origin "${BRANCH_NAME}:${BRANCH_NAME}_2" 2>&1)"
    git checkout -b "${BRANCH_NAME}_2" 2>/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit' >/dev/null
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    git checkout HEAD~2 >/dev/null 2>&1
    git checkout -b "${BRANCH_NAME}_3" 2>/dev/null
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit' >/dev/null
    echo "$FAKE_AWS_ACCESS_KEY_2" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    stdout="$(git push origin "${BRANCH_NAME}_2" "${BRANCH_NAME}_3" 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: detected secrets in branch" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_2_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing second plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "${BRANCH_NAME}_2" >/dev/null; then
        fail "$TEST_NAME —- missing first branch name in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "${BRANCH_NAME}_3" >/dev/null; then
        fail "$TEST_NAME —- missing second branch name in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep 'If the string is indeed sensitive' >/dev/null; then
        fail "$TEST_NAME —- missing remediation information in output for blocked push. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git checkout "$BRANCH_NAME" >/dev/null 2>&1
fi


TEST_NAME='should work if commits do not have any lines changed'
if true; then
    chmod +x "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit with no changes to content, just update file permissions' >/dev/null
    chmod -x "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit with no changes to content, just update file permissions' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code. Stdout and stderr:\n$stdout"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text"
    fi
    pass "$TEST_NAME"
fi

TEST_NAME='should work if commits do not add any new lines, only remove lines'
if true; then
    echo -e 'a\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm' > "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Ensure that file has multiple lines of content in it' >/dev/null
    git push --push-option=SKIP_LANGLOCK >/dev/null 2>&1
    grep -v 'b' "$NORMAL_FILE" > "$TEMP_FILE" && mv "$TEMP_FILE" "$NORMAL_FILE"
    git add "$NORMAL_FILE" && git commit -m 'Remove a single line from a file' >/dev/null
    grep -v 'f' "$NORMAL_FILE" > "$TEMP_FILE" && mv "$TEMP_FILE" "$NORMAL_FILE"
    git add "$NORMAL_FILE" && git commit -m 'Remove a single line from a file' >/dev/null
    grep -v 'j' "$NORMAL_FILE" > "$TEMP_FILE" && mv "$TEMP_FILE" "$NORMAL_FILE"
    git add "$NORMAL_FILE" && git commit -m 'Remove a single line from a file' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "remote: Server-side secret scanner (Langlock) detected no potential secrets in the push." >/dev/null; then
        fail "$TEST_NAME —- missing correct stdout text"
    fi
    pass "$TEST_NAME"
fi

TEST_NAME='should block push when detecting secret added in an intermediate commit but not present in final version of pushed code'
if true; then
    # Make a commit that adds a secret
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    # Make a commit that removes the same secret
    grep -v "$FAKE_AWS_ACCESS_KEY" "$NORMAL_FILE" > "$TEMP_FILE" && mv "$TEMP_FILE" "$NORMAL_FILE"
    git add "$NORMAL_FILE" && git commit -m 'Remove a single line from a file' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    if ! echo "$stdout" | grep "ERROR: detected secrets in branch" >/dev/null; then
        fail "$TEST_NAME —- missing summary text in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep "$FAKE_AWS_ACCESS_KEY_SHORT" >/dev/null; then
        fail "$TEST_NAME —- missing plaintext secret in output for blocked push. Instead found:\n$stdout"
    fi
    if ! echo "$stdout" | grep 'If the string is indeed sensitive' >/dev/null; then
        fail "$TEST_NAME —- missing remediation information in output for blocked push. Instead found:\n$stdout"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi

TEST_NAME='should only scan new commits, even if pushing a branch that got dangerous upstream commits merged in'
if true; then
    # Ensure that we have two branch, a regular one and a "secrets" one of which we will add secrets to, and then merge into the regular branch
    BRANCH_WITH_SECRETS_NAME="${BRANCH_NAME}_secrets"
    git checkout "$BRANCH_NAME" >/dev/null 2>&1
    removeAllowList
    git checkout -b "$BRANCH_WITH_SECRETS_NAME" >/dev/null 2>&1
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    git push -u origin "$BRANCH_WITH_SECRETS_NAME" --push-option=SKIP_LANGLOCK >/dev/null 2>&1
    git checkout "$BRANCH_NAME" >/dev/null 2>&1
    git merge "origin/${BRANCH_WITH_SECRETS_NAME}" 2>&1 >/dev/null || fail "Failed to setup test case"
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 0 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi


TEST_NAME='should scan the correct commits in currently pushed branch, even when there is a more recent merge from an upstream branch'
if true; then
    # Ensure that we have two branch, a regular one and a "secrets" one of which we will add secrets to, and then merge into the regular branch
    BRANCH_WITH_SECRETS_NAME="${BRANCH_NAME}_secrets_2"
    git checkout "$BRANCH_NAME" >/dev/null 2>&1
    removeAllowList
    git checkout -b "$BRANCH_WITH_SECRETS_NAME" >/dev/null 2>&1
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    git push -u origin "$BRANCH_WITH_SECRETS_NAME" --push-option=SKIP_LANGLOCK >/dev/null 2>&1
    git checkout "$BRANCH_NAME" >/dev/null 2>&1
    echo "$FAKE_AWS_ACCESS_KEY" >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make commit with credential' >/dev/null
    git merge "origin/${BRANCH_WITH_SECRETS_NAME}" 2>&1 >/dev/null || fail "Failed to setup test case"
    echo 'a' >> "$NORMAL_FILE" && git add "$NORMAL_FILE" && git commit -m 'Make benign commit' >/dev/null
    stdout="$(git push 2>&1)"
    if [ $? -ne 1 ]; then
        fail "$TEST_NAME —- wrong status code"
    fi
    pass "$TEST_NAME"
    # Cleanup
    git push --push-option=SKIP_LANGLOCK 2>/dev/null
fi

