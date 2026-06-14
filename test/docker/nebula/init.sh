#!/usr/bin/env bash
# init.sh — Start NebulaGraph with correct --local_ip and create test spaces.
# Usage: bash test/docker/nebula/init.sh
set -euo pipefail

TUTORIAL="$(cd "$(dirname "$0")/../../../tutorial" && pwd)"
export COMPOSE_FILE="$TUTORIAL/docker-compose.yml"

mkdir -p "$TUTORIAL/data/nebula/meta" \
        "$TUTORIAL/data/nebula/storage" \
        "$TUTORIAL/data/nebula/logs/meta" \
        "$TUTORIAL/data/nebula/logs/storage" \
        "$TUTORIAL/data/nebula/logs/graphd"

# Start metad first
echo "Starting metad..."
docker compose up -d nebula-metad
until docker inspect --format '{{.State.Health.Status}}' nebulagraph-metad 2>/dev/null | grep -q healthy; do sleep 2; done

# Now start storaged via docker-compose (will use wrong --local_ip)
echo "Starting storaged (temp)..."
docker compose up -d nebula-storaged
sleep 3

# Recreate storaged with correct --local_ip
STORAGED_IP=$(docker inspect nebulagraph-storaged --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
echo "Detected storaged IP: $STORAGED_IP"

docker compose stop nebula-storaged
docker compose rm -f nebula-storaged

docker run -d \
  --name nebulagraph-storaged \
  --network graphdb-net \
  -v "$TUTORIAL/data/nebula/storage:/data/storage" \
  -v "$TUTORIAL/data/nebula/logs/storage:/logs" \
  --user=root \
  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/vesoft/nebula-storaged:v3.8.0 \
  /usr/local/nebula/bin/nebula-storaged \
  --flagfile=/usr/local/nebula/etc/nebula-storaged.conf \
  --daemonize=false \
  --containerized=true \
  --meta_server_addrs=nebula-metad:9559 \
  --local_ip="$STORAGED_IP" \
  --ws_ip="$STORAGED_IP" \
  --port=9779 \
  --data_path=/data/storage \
  --log_dir=/logs \
  --v=0

echo "Waiting for storaged registration..."
sleep 15

# Start graphd
echo "Starting graphd..."
docker compose up -d nebulagraph
sleep 5

# ADD HOSTS and CREATE SPACE
echo "Configuring hosts and spaces..."
nebula-console -addr 127.0.0.1 -port 9669 -u root -p nebula \
  -e "ADD HOSTS $STORAGED_IP:9779" 2>&1
sleep 5

nebula-console -addr 127.0.0.1 -port 9669 -u root -p nebula \
  -e "CREATE SPACE IF NOT EXISTS graphx(partition_num=1, replica_factor=1, vid_type=FIXED_STRING(16))" 2>&1
nebula-console -addr 127.0.0.1 -port 9669 -u root -p nebula \
  -e "CREATE SPACE IF NOT EXISTS \`default\`(partition_num=1, replica_factor=1, vid_type=FIXED_STRING(16))" 2>&1

echo ""
echo "=== NebulaGraph Ready ==="
nebula-console -addr 127.0.0.1 -port 9669 -u root -p nebula -e "SHOW SPACES" 2>&1
