FROM alpine:3.12

RUN apk add --no-cache \
  ca-certificates \
  openssl \
  py3-pip \
  py3-setuptools 


COPY dnsrbl.py .
COPY requirements.txt .

RUN pip install -r requirements.txt

CMD [ "python3", "./dnsrbl.py" ]