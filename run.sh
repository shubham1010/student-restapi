
docker build -t student-restapi .
echo "student-restapi image is successfully created"
docker run -itd -p 8080:8080 --name student student-restapi
echo "Container is started"
docker exec -it student /bin/bash
