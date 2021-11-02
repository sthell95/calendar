#bin/sh
ssh-keygen -t rsa -b 4096 -m PEM -f ./config/auth/jwt.key -q -N ""