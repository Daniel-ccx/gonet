#!/bin/sh

##################################################
###   client proto & api
##################################################
awk -f proto.awk proto.txt > proto.go 
awk -f proto_func.awk proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk api.txt >> api.go 
awk -f api_rcode.awk api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet)([]byte, error){\n" >> api.go
awk -f api_bind_req.awk api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../agent/client_protos
mv -f api.go ../agent/client_protos

##################################################
### hub proto & api
##################################################
awk -f proto.awk hub_proto.txt > proto.go 
awk -f proto_func.awk hub_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk hub_api.txt >> api.go 
awk -f api_rcode.awk hub_api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(int32, *packet.Packet) ([]byte, error) = map[uint16]func(int32, *packet.Packet)([]byte, error){\n" >> api.go
awk -f api_bind_req.awk hub_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../hub/protos
mv -f api.go ../hub/protos

##################################################
### cooldown proto & api
##################################################
awk -f proto.awk cd_proto.txt > proto.go 
awk -f proto_func.awk cd_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk cd_api.txt >> api.go 
awk -f api_rcode.awk cd_api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*packet.Packet) ([]byte, error) = map[uint16]func(*packet.Packet)([]byte, error){\n" >> api.go
awk -f api_bind_req.awk cd_api.txt >> api.go 
printf "}" >> api.go

#### move #################
mv -f proto.go ../cooldown/protos
mv -f api.go ../cooldown/protos
