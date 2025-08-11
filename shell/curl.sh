# centos: ca证书默认目录
cat /etc/pki/tls/certs/ca-bundle.crt

# curl -v查看协议协议
curl -v https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation
curl -v https://api.stepfun.com/v1/chat/completions

# curl 时指定CA证书路径
curl --cacert /etc/pki/tls/certs/ca-bundle.crt.orig -v https://api.stepfun.com/v1/chat/completions
