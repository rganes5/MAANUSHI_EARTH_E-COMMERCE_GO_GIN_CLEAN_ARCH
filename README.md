# MAANUSHI_EARTH_E-COMMERCE_GO_GIN_CLEAN_ARCH

Back-end E-commerce REST API with Go & PostgresSQL. Uses Gin framework , Wire for dependency injection, Viper for handling environment variables, GORM as ORM and Swagger for API documentation.

Before using/cloning the folder, be sure to add a env file corresponding to the config.

# Technologies Used
Go
PostgreSQL
Gin
Wire
Viper
GORM
go-swagger

Clone the repo

git clone https://github.com/rganes5/MAANUSHI_EARTH_E-COMMERCE.git

# Install required packages

go mod tidy
# Setup Environment Variables

DB_HOST = replace with hostname
DB_NAME = replace with db name
DB_USER = replace with db username
DB_PORT = replace with db port
DB_PASSWORD = replace with db password

ACCOUNT_SID = replace with your twilio account sid
AUTHTOKEN = replace with twilio auth token
SERVICES_ID = replace with twilio services id
RAZORPAY_KEY= replace it with razorpay key
RAZORPAY_SECRET= replace it with razorpay secret
Compile and run

make run
