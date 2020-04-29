set -x

shift
cmd="./main"

until mysql -h127.0.0.1 -uroot -proot  -e 'source databaseConfig'; do
  >&2 echo "MySQL is unavailable - sleeping"
  sleep 1
done

>&2 echo "Mysql is up - executing command"
exec $cmd
