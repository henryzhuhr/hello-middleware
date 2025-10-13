#!/bin/bash

echo '等待 Redis 节点启动中...'
sleep 20

echo '检查所有节点是否可达...'
for node in redis-leader-0 redis-leader-1 redis-leader-2 redis-follower-0 redis-follower-1 redis-follower-2; do
    echo "检查节点 $node..."
    if ! redis-cli -h $node -p 6379 ping > /dev/null 2>&1; then
        echo "❌ 节点 $node 不可达"
        exit 1
    else
        echo "✅ 节点 $node 可达"
    fi
done

echo '正在创建 Redis Cluster...'
# 先重置所有节点的集群配置
for node in redis-leader-0 redis-leader-1 redis-leader-2 redis-follower-0 redis-follower-1 redis-follower-2; do
    redis-cli -h $node -p 6379 FLUSHALL > /dev/null 2>&1
    redis-cli -h $node -p 6379 CLUSTER RESET HARD > /dev/null 2>&1
done

sleep 5

# 使用 IP 地址创建集群
LEADER_0_IP=$(getent hosts redis-leader-0 | cut -d' ' -f1)
LEADER_1_IP=$(getent hosts redis-leader-1 | cut -d' ' -f1)
LEADER_2_IP=$(getent hosts redis-leader-2 | cut -d' ' -f1)
FOLLOWER_0_IP=$(getent hosts redis-follower-0 | cut -d' ' -f1)
FOLLOWER_1_IP=$(getent hosts redis-follower-1 | cut -d' ' -f1)
FOLLOWER_2_IP=$(getent hosts redis-follower-2 | cut -d' ' -f1)

echo "使用 IP 地址创建集群..."
redis-cli --cluster create \
    "${LEADER_0_IP}":6379 \
    "${LEADER_1_IP}":6379 \
    "${LEADER_2_IP}":6379 \
    "${FOLLOWER_0_IP}":6379 \
    "${FOLLOWER_1_IP}":6379 \
    "${FOLLOWER_2_IP}":6379 \
    --cluster-replicas 1 \
    --cluster-yes

echo '等待集群稳定...'
sleep 10

echo '验证集群状态...'
if redis-cli -h redis-leader-0 -p 6379 cluster info | grep -q 'cluster_state:ok'; then
    echo '✅ Redis Cluster 已成功初始化'
    echo '集群节点信息：'
    redis-cli -h redis-leader-0 -p 6379 cluster nodes
else
    echo '❌ Redis Cluster 初始化失败'
    echo '调试信息：'
    redis-cli -h redis-leader-0 -p 6379 cluster info
    redis-cli -h redis-leader-0 -p 6379 cluster nodes
    
    echo '尝试手动检查节点状态：'
    for node in redis-leader-0 redis-leader-1 redis-leader-2; do
        echo "=== $node 状态 ==="
        redis-cli -h $node -p 6379 cluster info | head -5
    done
fi

echo '保持容器运行...'
sleep infinity