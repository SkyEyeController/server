import re
from datetime import datetime

with open('data.out', 'r') as f:
    lines = f.readlines()

# 分别收集所有交易行和pre_node行
tx_lines = []
pre_node_lines = []

for line in lines:
    line = line.strip()
    if line.startswith("time_send:"):
        tx_lines.append(line)
    elif line.startswith("pre_node:"):
        pre_node_lines.append(line)

print(f"找到{len(tx_lines)}条交易记录和{len(pre_node_lines)}条pre_node记录")

# 分析每条交易记录
count = 0
threshold = 1  # 设置为1秒的阈值

# 处理每条交易记录
for i, tx_line in enumerate(tx_lines):
    # 提取rsend时间
    match_rsend = re.search(r'rsend:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})', tx_line)
    if not match_rsend:
        print(f"交易记录 #{i+1} 无法提取rsend时间: {tx_line}")
        continue
    
    rsend_time = match_rsend.group(1)
    rsend_dt = datetime.strptime(rsend_time, "%Y-%m-%d %H:%M:%S.%f")
    
    # 提取time_end时间，用于后续匹配
    match_end = re.search(r'time_end:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})', tx_line)
    if not match_end:
        print(f"交易记录 #{i+1} 无法提取time_end时间: {tx_line}")
        continue
    
    end_time = match_end.group(1)
    
    # 查找对应的pre_node行
    # 假设：与当前交易记录相关的pre_node行应该是time_end时间最接近的那条
    closest_pre_node = None
    min_time_diff = float('inf')
    
    for pre_node_line in pre_node_lines:
        # 解析pre_node行，提取最后一个节点的时间
        nodes = pre_node_line[9:].split(">")  # 去掉"pre_node:"前缀，按">"分割
        if not nodes:
            continue
        
        last_node = nodes[-1]
        node_parts = last_node.split("|")
        if len(node_parts) < 2:
            continue
        
        try:
            node_time = node_parts[1]
            node_dt = datetime.strptime(node_time, "%Y-%m-%d %H:%M:%S.%f")
            
            # 计算节点时间与交易end时间的接近程度
            time_diff = abs((datetime.strptime(end_time, "%Y-%m-%d %H:%M:%S.%f") - node_dt).total_seconds())
            
            if time_diff < min_time_diff:
                min_time_diff = time_diff
                closest_pre_node = pre_node_line
        except (ValueError, IndexError):
            continue
    
    if not closest_pre_node:
        print(f"交易记录 #{i+1} 找不到匹配的pre_node行: {tx_line}")
        continue
    
    # 现在我们有了匹配的pre_node行，提取最后一个节点的时间
    nodes = closest_pre_node[9:].split(">")
    last_node = nodes[-1]
    node_parts = last_node.split("|")
    
    try:
        node_time = node_parts[1]
        node_dt = datetime.strptime(node_time, "%Y-%m-%d %H:%M:%S.%f")
        
        # 计算rsend时间与最后节点时间的差值
        diff = abs((rsend_dt - node_dt).total_seconds())
        if diff > threshold:
            print(f"交易记录 #{i+1} rsend时间({rsend_time})与最后节点时间({node_time})差值超过{threshold}秒: {diff:.3f}秒")
            count += 1
    except (ValueError, IndexError):
        print(f"交易记录 #{i+1} 无法处理最后节点时间: {last_node}")
        continue

print(f"总共有{count}条记录的rsend与最后节点时间差超过{threshold}秒")