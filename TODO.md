# TODO
* Update CONTRIBUTING section in README
* Optimize regex performance 
* Add an allow list of config values that can be set on a per-repo basis
* Update readme with relevant git instructions
* Ensure that it fails safe in case of config issues, i.e., should still work if config is broken, maybe add validation
* Add test case for multiple users (clean up wording at top of test-e2e.sh)
* Finish e2e testing by allowing config options in the --push-option in the test mode
* langlock.sh command to scan all commits
* Maybe should have golang code accept multiple allow lists
* Maybe clean up gopath
* Test should detect secrets added in a merge
* Test should even if the current push did not touch the allow list, if the allow list has the wrong format, then the push will be blocked
* Test Should list all locations if the same secret shows up in multiple commits or files
* Test Should complain about misformatted allow list before secrets
* Test case for the global allow list
* Test case for testing the path-based one-liner
* Test should allow push if not enrolled users
* Test should block push if enrolled users but user not present
* Test should fail safe IF ENROLLED USERS HAS WRONG FORMAT
* Test should scan if user present in enrolled users with empty repos
* Test should scan if user present in enrolled users with repo in scanRepos
* Test should skip if user present in enrolled users with repo not in scanRepos
* Test should skip if user present in enrolled users with repo in skipRepos
* Test should scan if user present in enrolled users with repo not in skipRepos
* Test case for multiple repos
* Update docs for onboarding people one at a time