# 使用 awk 将日志中 all_t > 100 的日志全部捞取
awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); if(a[2] > 100) print $0}}'

# 按照日志中 all_t 的数值进行排序 +++是分隔符
cat * | head -100 | awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); print a[2],"+++ "$0}}' | sort -n

# 按照日志中 all_t 的数值进行排序 +++是分隔符 （不显示排序字段）
cat * | head -100 | awk -F "\t" '{for(i=1;i<=NF;i++) if($i ~ /^all_t=/){split($i, a, "="); print a[2],"+++ "$0}}' | sort -n | awk -F"+++" '{print $2}'