#!/bin/bash

echo "Waiting for RabbitMQ to be ready..."

until rabbitmqctl status; do
  echo "Still waiting..."
  sleep 2
done

echo "Creating exchange and queue..."
rabbitmqctl add_user admin 7866
rabbitmqctl set_user_tags admin administrator
rabbitmqctl delete_user guest

rabbitmqctl add_vhost barista
rabbitmqctl set_permissions -p barista admin ".*" ".*" ".*"

rabbitmqadmin declare exchange --vhost=barista name=gen_qr_event type=direct -u admin -p 7866 durable=true

rabbitmq-plugins enable rabbitmq_management

