# Edukasystem Backend Test
go run main.go <br>
hosted on localhost:8080 <br>
it uses 1 external libraries: <br>
github.com/julienschmidt/httprouter for the router


# API Calls
GET localhost:8080/MagicUpdate (part 1 of the specification) <br>
GET localhost:8080/SuperSorting (part 2 of the specification)


# How Magic Update Works
Magic update first generates a random number between 10000 and 20000, then it uses the check user api to check if the id is valid. If it isn't valid, it will generate the number again. Then if it is valid, the program would save the user ID and city ID in struct. After that it will generate a number between 0 and 1000 for the city ID and check if it is valid, if not it will generate again. After getting a valid city ID, the program will call the API in https://api.edukasystem.id/user/city as a seperate thread and serve the json response of user ID, City ID and Name before the PUT API call, and New City ID and Name after the API call. To get the city name, it uses the given API that returns the city name in json.

# How Super Sorting Works
The API first connects to the MySQL server, then runs a query to get the data and sorts it (uses SQL's ORDER BY), then parses the data to a slice of struct and then encodes it to json.

# How Logging Works
There is 2 functions that I made for logging, the first is just appending the user agent and access time to a log.txt file, and the other function is to clear the log.txt file every 15 seconds. These functions are called in a seperate thread (goroutines) to save time. Thankfully Golang can handle multithreading pretty easy (just add the go in front of the API call)

#Challenging Parts
I think that the most challenging part in making this is learning the new libraries and how to use them. Most of the concept that are used in this, I have already used in the past (but in a different language). So I had to google a lot just to get one thing implemented.