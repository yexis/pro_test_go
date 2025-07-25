# 使用 awk 将日志中 all_t > 100 的日志全部捞取
awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); if(a[2] > 100) print $0}}'

# 使用 awk解析日志，并解析固定字段trace_id
cat F_2025-07-25_10-54-48/* | awk -F"\t" '{for(i=1;i<=NF;i++) if($i ~ /^first_t=/){split($i, a, "="); if(a[2] > 3000) print $0}}' | grep -oP "\"trace_id\"\:\"\S+\"\,\"send" | awk -F "," '{print $1}'

# 按照日志中 all_t 的数值进行排序 +++是分隔符
cat * | head -100 | awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); print a[2],"+++ "$0}}' | sort -n

# 按照日志中 all_t 的数值进行排序 +++是分隔符 （不显示排序字段）
cat * | head -100 | awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); print a[2],"+++ "$0}}' | sort -n | awk -F"+++" '{print $2}'


# 计算 pv 比例
cat x.txt | grep -oP "first_t=\S+" | awk -F"=" '{
    total++
    if ($2 > 6000) count++
} END {
    if (total > 0)
        printf "比例: %.2f%%\n", count / total * 100
    else
        print "没有匹配项"
}'