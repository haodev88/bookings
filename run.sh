go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=postgres -cache=false -dbpass=root -production=false