after unpacking this project go to the dopc root
to launch the program run: go mod tidy && go run main.go

to use unit tests from all internal packages run: go test ./... 
unit tests were written with the help of generative AI

to test the whole program with curl for example:

distance less than 1km (one from assingment github)

curl "http://localhost:8080/api/v1/delivery-order-price?venue_slug=home-assignment-venue-helsinki&cart_value=1000&user_lat=60.17094&user_lon=24.93087"

distance above 1km 

curl "http://localhost:8080/api/v1/delivery-order-price?venue_slug=home-assignment-venue-helsinki&cart_value=1000&user_lat=60.1791&user_lon=24.9281"

distance above 1,5km

curl "http://localhost:8080/api/v1/delivery-order-price?venue_slug=home-assignment-venue-helsinki&cart_value=1000&user_lat=60.1840&user_lon=24.9281"

distance above 2km

curl "http://localhost:8080/api/v1/delivery-order-price?venue_slug=home-assignment-venue-helsinki&cart_value=1000&user_lat=60.1901&user_lon=24.9281"