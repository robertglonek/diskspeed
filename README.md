# Read disk one block at a time and report speed at intervals

## Usage

```
% ./diskspeed 
ERROR: Usage: diskspeed [optional-params] /dev/disk
  -print-sleep int
        seconds between speed sampling and prints (default 1)
  -read-size int
        size in bytes of read at a time (default 1048576)
```

Works with both disks and files (use file path instead of `/dev/disk` to check file read speed)

## Examples

```
% ./diskspeed testfile.txt
1840 MB/s for 1 s
1712 MB/s for 1 s
1697 MB/s for 1 s
1764 MB/s for 1 s
1824 MB/s for 1 s
1765 MB/s for 1 s
503 MB/s for 1 s
END

% ./diskspeed /dev/sdx
1661 MB/s for 1 s
1632 MB/s for 1 s
1655 MB/s for 1 s
1529 MB/s for 1 s
1568 MB/s for 1 s
1738 MB/s for 1 s
1322 MB/s for 1 s
END
```
