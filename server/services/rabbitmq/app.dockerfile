FROM rabbitmq:3.11-management

COPY setup.sh /setup.sh

RUN chmod +x /setup.sh

CMD ["sh", "-c", "rabbitmq-server & sleep 10 && /setup.sh && wait"]