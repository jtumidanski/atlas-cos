if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache -f Dockerfile.dev --tag atlas-cos:latest .
else
   docker build -f Dockerfile.dev --tag atlas-cos:latest .
fi
