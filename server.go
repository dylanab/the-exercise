package main;

import (
	"fmt"
	"log"
	"github.com/valyala/fasthttp"
);

const discovrHost string = "https://epic.gsfc.nasa.gov";
const discovrAvailableDates string = "/api/natural/all"
const discovrDate string = discovrHost + "/api/natural/YYYY/MM/DD";

//TODO(Dylan): replace this with an actual cache 
//var discovrAvailableDatesCache []byte;
//var discovrLatestDateCache []byte;
var renderedHTMLCache []byte;

//Note(Dylan): probably break these processes out into different functions and/or different packages
//Note(Dylan): this can be a goroutine that gets called periodically. Though we do actually want this to block in the case that the cache is empty when we get a GET request.
func FetchRenderAndCache() {

	//TODO(Dylan): request and cache discovrAvailableDates. The data will change once per day, so we can hold it in the memcache for a really long time.
	
	//TODO(Dylan): request and cache the detailed discovrData for the latest date

	//TODO(Dylan): parse the discovrData into some go structs that will help us render the HTML

	//TODO(Dylan): render some html and store it in the cache
}

func main() {
	
	//NOTE(Dylan): initialize the cache
	renderedHTMLCache := []byte("");

	//TODO(Dylan): Before the server starts up, request the discovrLatest data and fill the cache.
	FetchRenderAndCache();

	server := func(ctx *fasthttp.RequestCtx) {
		//TODO(Dylan): if the cache is empty or expired, preform the FetchRenderandCache procedure synchronously
		ctx.Response.Header.Set("Content-Type", "text/plain");
		ctx.Response.SetBody(renderedHTMLCache);
	}

	err := fasthttp.ListenAndServe(":8080", server);
	exitMessage := "exercise server has stopped";
	if err != nil {
		exitMessage = fmt.Sprintf("%s with error: %s", exitMessage, err);
	}
	log.Fatalf("%s", exitMessage);
}

