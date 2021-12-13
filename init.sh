rm -r *.db 
docker container stop $(docker container ls -q)
docker container prune