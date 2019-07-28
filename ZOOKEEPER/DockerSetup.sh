
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi -f $(docker images -q)

docker run --name some-zookeeper --network host --restart always -d zookeeper

#docker build --file DockerFile.ZooKeeper --tag zookeeper_image .

#docker create --network host --name zookeeper_server_1 zookeeper_image
#docker create --network host --name zookeeper_server_2 zookeeper_image
#docker create --network host --name zookeeper_server_3 zookeeper_image

#docker container start zookeeper_server_1
#docker container start zookeeper_server_2
#docker container start zookeeper_server_3


# docker exec zookeeper_server_2 /bin/zkServer.sh start
# docker exec -it some-zookeeper /bin/bash
# docker exec -it primary1 /go/primarynode
# zkCli.sh


#docker stop zookeeper_server_3


docker build --file DockerFile.PrimaryNode --tag SDCC_Project_PrimaryNode_image .
