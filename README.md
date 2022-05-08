# :lock_with_ink_pen: Entities

Modern go identity manager for UNIX systems.

Entities parses includes file to generate UNIX-compliant `/etc/passwd` , `/etc/shadow` and `/etc/groups` files.
It can be used to handle identities management and honors already existing entities in the system.


```

$> entities apply <entity.yaml>
$> entities delete <entity.yaml>
$> entities create <entity.yaml>

```

## Entities file format

### Passwd

```yaml
kind: "user"
username: "foo"
password: "pass"
uid: 0
gid: 0
info: "Foo!"
homedir: "/home/foo"
shell: "/bin/bash"
```

To use dynamic uid allocation set the `uid` field with value `-1`:

```yaml
kind: "user"
username: "foo"
password: "pass"
uid: -1
gid: 500
info: "Foo!"
homedir: "/home/foo"
shell: "/bin/bash"
```

`entities` will searching for the first available range specified by the env variable
`ENTITY_DYNAMIC_RANGE` or by the default the range `500-999`.


To set gid with a dynamic id based by the group name you can set the `group` attribute:
```yaml
kind: "user"
username: "foo"
password: "pass"
uid: 100
group: "foogroup"
info: "Foo!"
homedir: "/home/foo"
shell: "/bin/bash"
```

`entities` will retrieve the `gid` from existing `/etc/group` file.


### Gshadow

```yaml
kind: "gshadow"
name: "postmaster"
password: "foo"
administrators: "barred"
members: "baz"
```

### Shadow

```yaml
kind: "shadow"
username: "foo"
password: "bar"
last_changed: 1
minimum_changed: 2
maximum_changed: 3
warn: 4
inactive: 5
expire: 6
```

To define `last_changed` with a value equal to current days from 1970 use `now`.

### Group

```yaml
kind: "group"
group_name: "sddm"
password: "xx"
gid: 1
users: "one,two,tree"
```

To assign a dynamic gid it's possible to use the value `-1`:

```yaml
kind: "group"
group_name: "foogroup"
password: "xx"
gid: -1
users: "one,two,tree"
```

`entities` will searching for the first available range specified by the env variable
`ENTITY_DYNAMIC_RANGE` or by the default the range `500-999`.

### List entities

To read and list entities available in a system (users, groups, shadow, gshadow):

```shell
$> entities list users
+----------------+--------------------+---------+----------+--------------------------------+--------------------------+----------------+
|    USERNAME    | ENCRYPTED PASSWORD | USER ID | GROUP ID |              INFO              |         HOMEDIR          |     SHELL      |
+----------------+--------------------+---------+----------+--------------------------------+--------------------------+----------------+
| adm            | x                  |       3 |        4 | adm                            | /var/adm                 | /bin/false     |
| apache         | x                  |      81 |       81 | added by portage for apache    | /var/www                 | /sbin/nologin  |
| arangodb3      | x                  |    1001 |     1006 |                                | /home/arangodb3          | /bin/bash      |
| avahi          | x                  |     104 |      104 | added by portage for avahi     | /dev/null                | /sbin/nologin  |
| bin            | x                  |       1 |        1 | bin                            | /bin                     | /bin/false     |
...
```

`entities` by default read files `/etc/passwd`, `/etc/groups`, `/etc/gshadow` and `/etc/shadow.

To read entities from a different file use `-f|--file`:

```shell
$> entities list users --file /tmp/passwd

$> # Read list of available groups
$> entities list groups

$> # Read list of available groups order by id
$> entities list groups -s id

$> # Read list of available groups order by name
$> entities list groups -s name

$> # Read list of gshadow entries
$> entities list gshadow

$> # Read list of shadow entries
$> entities list shadow
```

`entities` permits to list entities defined in YAML from a directory too:

```shell
$> entities list users --specs-dir /entities-catalog
```

### Dump entities

`entities` permits to generate `entities` specs from existing rootfs:

```shell
$> entities dump -t ./catalog
Creating 41 users under the directory catalog/users
Creating 70 groups under the directory catalog/groups
Creating 41 shadows under the directory catalog/shadows
Creating 13 gshadows under the directory catalog/gshadows
All done.
```

or from specified files:

```shell
$> entities dump -t ./catalog --groups-file /tmp/groups --gshadow-file /tmp/gshadow --shadow-file /tmp/shadow --users-file /tmp/passwd

### Merge entities

The idea of the `merge` subcommand is to use an existing **catalog** and then merge entities if they aren't yet present.

```shell
$> # merge all entities defined on a catalog on /etc/passwd,/etc/groups,/etc/shadow,/etc/gshadow
$> entities merge --specs-dir ./my-catalog -a

$> # merge all entities defined on a catalog on custom files
$> entities merge --specs-dir ./my-catalog -a --groups-file /tmp/groups --gshadow-file /tmp/gshadow --users-file /tmp/passwd --shadow-file /tmp/shadow

$> # merge all entry related with a specific entity defined on a catalog on /etc/passwd,/etc/groups,/etc/shadow,/etc/gshadow.
$> # On the example is created the group mongodb
$> entities merge --specs-dir ./my-catalog -e mongodb
```
