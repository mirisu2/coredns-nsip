# coredns-nsip
Corefile
```
foo.bar:1053 {
    log
    nsip {
        ns1 192.168.1.100
        ns2 192.168.2.101
    }
}
foo.baz:1053 {
    log
    nsip {
        ns1 192.168.2.102
    }
    errors
}
```
kubernetes
```
$ kubectl -n ns1 exec -it nginx -- bash
root@nginx:/# curl -v a.foo.bar
*   Trying 192.168.1.100:80...

$ kubectl -n ns2 exec -it nginx -- bash
root@nginx:/# curl -v a.foo.bar
*   Trying 192.168.1.101:80...
```
