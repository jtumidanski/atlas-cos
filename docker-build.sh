if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache --tag atlas-cos:latest .
else
   docker build --tag atlas-cos:latest .
fi
