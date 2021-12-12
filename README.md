# gobackup
遍历给定父目录下的文件，压缩后并删除原文件

## 简单用法

遍历`/app/test`目录下的文件, 并设置协程为2（同时压缩2个）

```
gobackup -path /app/test -t 2
```

## 命令行选项

```
-path string
        Traverse the directory under the path compression path, delete after compression.
-t    int
        Set number of concurrent coroutines. (default 2)
```
