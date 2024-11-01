This is a RSS aggregator CLI tool that stores subscribed feeds into a sql database. 
To run this program you will need to have installed postgres and Go. 
To install the gator CLI tool run the command go install https://github.com/ecastellanosr/rssagg/
To use the database you need a .gatorconfig.json file containing the database url and the current user. Ej. 
{"db_url":"postgres://user:password@localhost:5432/gator?sslmode=disable","current_user_name":"name"}
Some commands you can run in the CLI tool are: 
1. Register: 
2. Login:
3. addfeed:
4. 
5.
6.
