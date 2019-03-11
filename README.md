# ansible_connection_test
Test http/network connectivity from local and remote linux hosts using Ansible

These modules are written in `go` for the purpose of portability. This is nice for hosts that do not have Python installed.

Ansible will push the binaries to remote hosts if you drop them in a `library` directory in your playbook path.

# Example
```yaml
---
- hosts: localhost
  name: runs go binary on macos host
  tasks:
    - check_http_osx:
        checks: &checks
          - name: google.ca/
            url: https://google.ca
            no_follow_redirect: true
            expected:
              status_code: 301

- hosts: remote_linux_host
  name: runs go binary on linux host
  tasks:
    - check_http_linux:
        checks: *check
```

```sh
ansible-playbook test.yml
```

