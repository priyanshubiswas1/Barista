#!/bin/bash

# Configuration from environment variables with defaults
USER=${RABBITMQ_USER:-admin}
PASS=${RABBITMQ_PASS:-7866}
VHOST=${RABBITMQ_VHOST:-barista}

echo "Waiting for RabbitMQ to be fully started..."

# Robust wait loop using rabbitmq-diagnostics
MAX_RETRIES=30
RETRY_COUNT=0
until rabbitmq-diagnostics -q check_running; do
  RETRY_COUNT=$((RETRY_COUNT+1))
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo "Error: RabbitMQ failed to start after $MAX_RETRIES attempts."
    exit 1
  fi
  echo "Still waiting ($RETRY_COUNT/$MAX_RETRIES)..."
  sleep 2
done

# Small extra buffer to ensure management plugin and internal apps are ready
sleep 5

echo "Configuring RabbitMQ for user: $USER and vhost: $VHOST..."

# Create user and set permissions
rabbitmqctl add_user $USER $PASS || echo "User $USER already exists"
rabbitmqctl set_user_tags $USER administrator

# Remove default guest user
rabbitmqctl delete_user guest || echo "Guest user already removed"

# Create vhost and set permissions
rabbitmqctl add_vhost $VHOST || echo "Vhost $VHOST already exists"
rabbitmqctl set_permissions -p $VHOST $USER ".*" ".*" ".*"

# Declare exchange using rabbitmqadmin
rabbitmqadmin declare exchange --vhost=$VHOST name=gen_qr_event type=direct -u $USER -p $PASS durable=true

# Enable management plugin (usually enabled in the base image, but good to ensure)
rabbitmq-plugins enable rabbitmq_management

