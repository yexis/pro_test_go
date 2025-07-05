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

##### 删除本地 tag，同时删除远端
```bash
git push origin --delete v1.0.0
```