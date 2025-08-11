# 免密登录机器，把公钥拷到服务端的~/.ssh/authorized_keys里面就好

# 1. 本地生成ssh公钥+秘钥
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
cat ~/.ssh/id_rsa.pub

# 2. 服务端新增识别文件
~/.ssh/authorized_keys

# 3. 将本地公钥写到远端识别文件
vim ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# 4. 登录
ssh work@....