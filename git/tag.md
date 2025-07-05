##### 创建一个 tag
```bash
git tag v1.0.0
```

##### 将 tag 推送到远端
```bash
git push origin v1.0.0
```

##### 将所有 tag 推送到远端
```bash
git push origin --tags
```

##### 删除本地 tag，不影响远端
```bash
git tag -d v1.0.0
```

##### 删除远端 tag
```bash
git push origin --delete v1.0.0
```


##### 给子模块打 tag （git tag 的时候带上路径前缀）
比如说 pro_test_go 和 pro_test_go/easy 同时作为 module 发布了；
用户可以直接依赖 pro_test_go，也可以直接依赖 pro_test_go/easy；
```bash
git tag easy/v1.0.1
git push origin easy/v1.0.1
```
然后在业务模块中执行
```bash
go get github.com/yexis/pro_test_go/easy@v1.0.1
```
