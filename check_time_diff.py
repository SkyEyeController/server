import re
from datetime import datetime

with open('data.out', 'r') as f:
    lines = f.readlines()

count = 0
mode = ">"
threshold = 2

for i in range(0, len(lines), 2):
    line = lines[i]
    match_rsend = re.search(r'rsend:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})', line)
    match_end = re.search(r'time_end:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})', line)
    if match_rsend and match_end:
        rsend_dt = datetime.strptime(match_rsend.group(1), "%Y-%m-%d %H:%M:%S.%f")
        end_dt = datetime.strptime(match_end.group(1), "%Y-%m-%d %H:%M:%S.%f")
        diff = (end_dt - rsend_dt).total_seconds()
        if mode == ">":
            if diff > threshold:
                print(f"第{i//2+1}条记录 rsend 和 time_end 时间差超过{threshold}秒: {diff:.3f}秒")
                count += 1
        else:
            if diff < threshold:
                print(f"第{i//2+1}条记录 rsend 和 time_end 时间差超过{threshold}秒: {diff:.3f}秒")
                count += 1

if mode == ">":
    print(f"总共有{count}条记录 rsend 和 time_end 时间差超过{threshold}秒")
else :
    print(f"总共有{count}条记录 rsend 和 time_end 时间差小于{threshold}秒")




# TranRate
# 修改export
# 修改逻辑