cd frontend
/root/.deno/bin/deno install
/root/.deno/bin/deno task build --emptyOutDir
cd ..
/usr/local/go/bin/go run main.go
