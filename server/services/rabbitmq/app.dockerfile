FROM rabbitmq:3.11-management

COPY setup.sh /setup.sh

# Ensure script is executable and has Unix line endings
RUN chmod +x /setup.sh \
    && sed -i 's/\r$//' /setup.sh

CMD ["sh", "-c", "rabbitmq-server & sleep 10 && echo 'Listing root dir:' && ls -la / && echo 'Attempting to run /setup.sh' && sh -x /setup.sh && wait"]