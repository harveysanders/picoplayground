FROM eclipse-mosquitto:openssl

RUN mkdir /etc/mosquitto
ADD mosquitto.conf /etc/mosquitto
ADD acl.conf /etc/mosquitto/acl.conf
ADD passwd /etc/mosquitto

RUN chmod 0700 /etc/mosquitto/acl.conf && chown mosquitto:mosquitto /etc/mosquitto/acl.conf
RUN chmod 640 /etc/mosquitto/passwd && chown mosquitto:mosquitto /etc/mosquitto/passwd

CMD ["/usr/sbin/mosquitto", "-c", "/etc/mosquitto/mosquitto.conf"]
