# 查看ca证书chain
openssl s_client -connect api.stepfun.com:443 -showcerts | less