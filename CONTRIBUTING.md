# Langlock — (DEPRECATED) CONTRIBUTING

The instructions for contributing need to be updated.

## Overview
Langlock is a server-side pre-receive hook that prevents developers from accidentally pushing hardcoded secrets upstream. We named it after a jinx that Harry Potter uses to glue the target’s tongue to the roof of their mouth.

Langlock is currently in development. We have not deployed it to any repos yet.

## Related documentation

- Please reach out to the authors for access to architectural diagrams and sequence diagrams.

## Getting started (running the pre-receive hook locally)

Follow the instructions [here](https://help.github.com/en/enterprise/2.19/admin/developer-workflow/creating-a-pre-receive-hook-script#testing-pre-receive-scripts-locally), but make the following modifications:

- use the `Dockerfile` in this repo, instead of their `Dockerfile.dev`
- use the `pre-receive-hook.sh` file in this repo, instead of their `always_reject.sh` file
- note that the `pre-receive-hook.sh` file in this repo contains a placeholder string ("PUT_SLACK_WEBHOOK_URL_HERE") where the slack webhook url is needed. Rather than hardcoding that URL into the repo, we can provide the url dynamically when we upload the script to the docker volume that runs the pre-receive hook, as shown (note that docker cp does not support bash's process substitution, so I use an auxiliary file instead):

        SLACK_WEBHOOK_URL=...
        
        cat pre-receive-hook.sh | sed "s|PUT_SLACK_WEBHOOK_URL_HERE|$SLACK_WEBHOOK_URL|g" > temp.pre-receive-hook.sh; chmod +x temp.pre-receive-hook.sh; docker cp temp.pre-receive-hook.sh data:/home/git/test.git/hooks/pre-receive

    You can find the specific webhook url by visiting the [Slack web console](https://api.slack.com/apps/SLACKAPPID/incoming-webhooks).

- when running the provided push command, make sure to add the CLI option `--push-option=SECRET_SCAN`, in order to opt-in to the secret scanning logic. Eg

        GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p 52311 -i ../id_rsa" git push --push-option=SECRET_SCAN -u test master:master

- ask to join the Slack channel #test-secret-scanner-alerts to see the output of the webhook during development

**Troubleshooting**

- Ensure that the `pre-receive-hook.sh` script that we upload to the git server is executable. ie.

        docker container exec -it NAME_OF_YOUR_GIT_SERVER_CONTAINER /bin/bash
        ls -la test.git/hooks/pre-receive
        chmod +x test.git/hooks/pre-receive
        exit

- Ensure that the git server has pre-receive hooks enabled for the repo. The `Dockerfile` should take care of this automatically.

        docker container exec -it NAME_OF_YOUR_GIT_SERVER_CONTAINER /bin/bash
        cd test.git
        git config receive.advertisePushOptions true
        exit

## Useful aliases during local development

    # Push, with secret scanning enabled
    alias p='GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p 52311 -i ../id_rsa" git push --push-option=SECRET_SCAN -u test master:master'
    
    # Copy our version of pre-receive hook to the git server container
    SLACK_WEBHOOK_URL=...
    alias c='cat pre-receive-hook.sh | sed "s|PUT_SLACK_WEBHOOK_URL_HERE|$SLACK_WEBHOOK_URL|g" > temp.pre-receive-hook.sh; chmod +x temp.pre-receive-hook.sh; docker cp temp.pre-receive-hook.sh data:/home/git/test.git/hooks/pre-receive'
    
    # Make a non-sensitive commit in the demo repo
    alias e='echo "Benign change" >> README; git add -u; git commit -m "Make non-sensitive change";'
    
    # Make a potentially sensitive commit in the demo repo
    alias es='echo "Foo CONFIDENTIAL bar" >> SOMEFILE; git add -u; git commit -m "Make sensitive change potentially including a secret";'
    
    # Make a commit to the the allow-list file in the demo repo
    alias ea='echo "New secret" >> secrets-allow-list.json; git add -u; git commit -m "Add secret to allow list";'
    
    # Get an interactive shell in the git server container
    alias i='docker container exec -it NAME_OF_YOUR_GIT_SERVER_CONTAINER /bin/sh'
    
    # Run the git server container and give it a specific name
    alias r='docker run -d -p 52311:22 --volumes-from data --name DESIRED_NAME_OF_YOUR_GIT_SERVER_CONTAINER pre-receive.dev'
