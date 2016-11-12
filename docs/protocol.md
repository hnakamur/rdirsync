
## readDir

request

```
readDir\t<path>
```

response

<size> is empty for directories

```
count\t<entryCount>
entry\t<mode>\t<atime>\t<mtime>\t<owner>\t<group>\t<size>\t<name>
entry\t<mode>\t<atime>\t<mtime>\t<owner>\t<group>\t<size>\t<name>
entry\t<mode>\t<atime>\t<mtime>\t<owner>\t<group>\t\t<name>
...
```

```
err\t<message>
```

## sendFile

request

```
sendFile\t<mode>\t<atime>\t<mtime>\t<owner>\t<group>\t<size>\t<path>
<data>
```

response

```
ok
```

```
err\t<message>
```

## remove

```
remove\t<path>
```

response

```
ok
```

```
err\t<message>
```

## sendDir

request

```
sendDir\t<mode>\t<atime>\t<mtime>\t<owner>\t<group>\t<path>
```

response

```
ok
```

```
err\t<message>
```
