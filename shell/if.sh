
# 判断变量是否为空 if [ -z "$sear" ];
cat x9_pro_lite_a.txt | while read -r line; do
  sear=$(grep $line cuid_x9prolite_rtc.txt);
  if [ -z "$sear" ]; then
    echo $sear;
  fi;
done

# 判断变量是否非空 if [ -n "$sear" ];
cat x9_pro_lite_a.txt | while read -r line; do
  sear=$(grep $line cuid_x9prolite_rtc.txt);
  if [ -n "$sear" ]; then
    echo $sear;
  fi;
done
