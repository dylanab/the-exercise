# the-exercise

Instructions:

1) build the docker image
  $ docker build -t the-exercise . 

2) run the docker image, mapping localhost:8080 to port 8080 in the container
  $ docker run -p 8080:8080 the-exercise

3) Use a browser to navigate to localhost:8080/ 
   You should see a page of date links. Each of these links routes to a page showing the DISCOVR Blue Marble images for the given date. You can also navigate directly to a date page with the following url pattern: localhost:8080/date/YYYY-MM-DD
