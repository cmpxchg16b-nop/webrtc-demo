cfssl genkey ca.json | cfssljson -bare ca
cfssl genkey -initca ca.json | cfssljson -bare ca
openssl x509 -in ca.pem -noout -text
