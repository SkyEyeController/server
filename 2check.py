import re
from datetime import datetime

with open('data.out', 'r') as f:
    lines = f.readlines()

count = 0
threshold = 1  # 设置为1秒的阈值

# 确保处理的是完整记录
if len(lines) % 2 != 0:
    print(f"警告：数据文件行数不是偶数 ({len(lines)}行)")

for i in range(0, len(lines)-1, 2):
    tx_line = lines[i]
    pre_node_line = lines[i+1]
    
    # 提取rsend时间
    match_rsend = re.search(r'rsend:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})', tx_line)
    if not match_rsend:
        print(f"第{i//2+1}条记录无法提取rsend时间: {tx_line.strip()}")
        continue
    
    rsend_time = match_rsend.group(1)
    rsend_dt = datetime.strptime(rsend_time, "%Y-%m-%d %H:%M:%S.%f")
    
    # 确认下一行是pre_node行
    if not pre_node_line.strip().startswith("pre_node:"):
        print(f"第{i//2+1}条记录的下一行不是pre_node行: {pre_node_line.strip()}")
        continue
    
    # 解析pre_node行，提取最后一个节点的时间
    nodes = pre_node_line.strip()[9:].split(">")  # 去掉"pre_node:"前缀，按">"分割
    if not nodes:
        print(f"第{i//2+1}条记录的pre_node行没有节点信息: {pre_node_line.strip()}")
        continue
    
    last_node = nodes[-1]
    node_parts = last_node.split("|")
    if len(node_parts) < 2:
        print(f"第{i//2+1}条记录的最后节点格式不正确: {last_node}")
        continue
    
    node_time = node_parts[1]
    try:
        node_dt = datetime.strptime(node_time, "%Y-%m-%d %H:%M:%S.%f")
    except ValueError:
        print(f"第{i//2+1}条记录的节点时间格式不正确: {node_time}")
        continue
    
    # 计算时间差
    diff = abs((rsend_dt - node_dt).total_seconds())
    if diff > threshold:
        print(f"第{i//2+1}条记录 rsend时间({rsend_time})与最后节点时间({node_time})差值超过{threshold}秒: {diff:.3f}秒")
        count += 1

print(f"总共有{count}条记录的rsend与最后节点时间差超过{threshold}秒")