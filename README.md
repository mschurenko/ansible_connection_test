# ansible_connection_test
test http/network connectivity from linux hosts using ansible/go

# Usage
```sh
$ ansible-playbook -i 'localhost,bastion,' -v main.yml
Using /Users/mschurenko/git/mschurenko/ansible_connection_test/ansible.cfg as config file

PLAY [localhost] ********************************************************************************************************

TASK [check_http_osx] ***************************************************************************************************
ok: [localhost] => changed=false
  checks:
    /: |-
      status code: 301 matches 301
    /foo: |-
      status code: 404 matches 404

PLAY [bastion] **********************************************************************************************************

TASK [check_http] *******************************************************************************************************
ok: [bastion] => changed=false
  checks:
    /: |-
      status code: 301 matches 301
    /foo: |-
      status code: 404 matches 404

PLAY RECAP **************************************************************************************************************
bastion                    : ok=1    changed=0    unreachable=0    failed=0
localhost                  : ok=1    changed=0    unreachable=0    failed=0
```
