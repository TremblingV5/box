# Example for configx with etcd

1. Use `docker-compose.yml` to start an etcd server
2. Set environment variable
   - `CONFIG_TYPE=ETCD`
   - `CONFIG_CENTER_URL=http://localhost:2379`
3. Open `config/config.json` and copy the content to etcd key `CONFIG_CENTER_STORE_KEY_PREFIX/config`
4. Open `config/biz.json` and copy the content to etcd key `CONFIG_CENTER_STORE_KEY_PREFIX/biz`
5. Run the example, it will output the same value of the 2 json files
