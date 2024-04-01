# coredns-nsip
Corefile
```
foo.bar:1053 {
    loop
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
Coredns rebuild
```
$ git clone https://github.com/coredns/coredns.git
$ cd coredns
# configure Corefile
$ vim Corefile

$ cat plugin.cfg
nsip:github.com/mirisu2/coredns-nsip
forward:forward

$ GOPROXY=direct go get github.com/mirisu2/coredns-nsip
$ make
$ docker build --network host -t dockeruser/coredns:nsip .
$ docker push dockeruser/coredns:nsip
```
