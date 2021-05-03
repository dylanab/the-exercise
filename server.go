package main;

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"html/template"
	"github.com/valyala/fasthttp"
	"exercise/httpclient"
	"exercise/cache"
);

const port string = ":8080";
const discovrHost string = "https://epic.gsfc.nasa.gov";
const discovrAvailableDates string = discovrHost + "/api/natural/all"
const discovrDate string = discovrHost + "/api/natural/date/";

var c *cache.Cache;
var cErr error;

var availabilityTemplate *template.Template;
var atErr error;
var datePageTemplate *template.Template;
var dtErr error;

type DiscovrAvailability struct {
    Date string `json:"date"`;
}

type DiscovrImage struct {
	Url string `json:"image"`;
	Timestamp string `json:"date"`;
	YYYY, MM, DD string;
}

func RenderAvailabilityPage(ctx *fasthttp.RequestCtx, c *cache.Cache) {
	//NOTE(Dylan): check if the availiability html is in the cache, if it hasn't or it's expired, then we need to pull the data from the DISCOVR api and render it into HTML. 
	var discovrAvailability []DiscovrAvailability;
	var getErr error;
	availableDatesBytes, isExpired := c.GetItem("availability");
	
	if(isExpired) {
		log.Printf("getting availableDatesBytes");
		availableDatesBytes, _, _, getErr = httpclient.GetBytes(discovrAvailableDates, 3000);
		if(getErr != nil) {
			log.Printf("%s", getErr);
			return;	
		}	
		c.AddItem("availability", availableDatesBytes);
	}

	
	//log.Printf("%s", availableDatesBytes);
	jErr := json.Unmarshal(availableDatesBytes, &discovrAvailability);

	if(jErr != nil) {		
		jErrMsg := fmt.Sprintf("Error during availability page json unmarshal %s", jErr);
		log.Print(jErrMsg);
		ctx.Error(jErrMsg, fasthttp.StatusInternalServerError);
		return;
	}

	//log.Printf("unmarshalled availability json: %+v", discovrAvailability);

	log.Printf("attempting to execute availabilitypage template");
	tErr := availabilityTemplate.ExecuteTemplate(ctx, "availabilitypage", discovrAvailability);

	if(tErr != nil) {
		tErrMsg := fmt.Sprintf("Error during availability page template render %s", tErr);
		log.Print(tErrMsg);
		ctx.Error(tErrMsg, fasthttp.StatusInternalServerError);
		return;
	}

	ctx.SetContentType("text/html");
	ctx.SetStatusCode(fasthttp.StatusOK);
}

func RenderDatePage(ctx *fasthttp.RequestCtx, datestring string, c *cache.Cache) {

	var discovrDateImages []DiscovrImage;
	var getErr error;
	dateBytes, isExpired := c.GetItem(datestring);
	
	var YYYY string = "";
	var MM string = "";
	var DD string = "";

	dateStringSplit := strings.Split(datestring, "-");
	if(len(dateStringSplit) == 3) {
		YYYY = dateStringSplit[0];
		MM = dateStringSplit[1];
		DD = dateStringSplit[2];
	}
	
	log.Printf("YYYY: %s, MM: %s, DD: %s", YYYY, MM, DD);

	if(isExpired) {
		dateBytes, _, _, getErr = httpclient.GetBytes(discovrDate + datestring, 3000);
		if(getErr != nil) {
			log.Printf("%s", getErr);
			return;	
		}	
		c.AddItem(datestring, dateBytes);
	}

	//log.Printf("%s", dateBytes);
	jErr := json.Unmarshal(dateBytes, &discovrDateImages);

	for i := range(discovrDateImages) {
		image := &discovrDateImages[i];
		image.YYYY = YYYY;
		image.MM = MM;
		image.DD = DD;
	}

	log.Printf("%+v", discovrDateImages[0]);

	if(jErr != nil) {		
		jErrMsg := fmt.Sprintf("Error during date page json unmarshal %s", jErr);
		log.Print(jErrMsg);
		ctx.Error(jErrMsg, fasthttp.StatusInternalServerError);
		return;
	}

	tErr := datePageTemplate.ExecuteTemplate(ctx, "datepage", discovrDateImages);
	
	if(tErr != nil) {
		tErrMsg := fmt.Sprintf("Error during date page template render %s", tErr);
		log.Print(tErrMsg);
		ctx.Error(tErrMsg, fasthttp.StatusInternalServerError);
		return;
	}
    
	ctx.SetContentType("text/html");
	ctx.SetStatusCode(fasthttp.StatusOK);
}


func main() {

	startErrMessage := "the exercise server failed to start";

	c, cErr = cache.Create("120s", true);
	
	if(cErr != nil) {
		log.Fatalf(startErrMessage + " with error  %s", cErr);
	}

	availabilityTemplate, atErr = template.ParseFiles("templates/availabledates.gohtml");
	datePageTemplate, dtErr = template.ParseFiles("templates/datepage.gohtml");

	if(atErr != nil) {
		log.Fatalf(startErrMessage + " with error  %s", atErr);
	}

	if(dtErr != nil) {
		log.Fatalf(startErrMessage + " with error  %s", dtErr);
	}
		
	server := func(ctx *fasthttp.RequestCtx) {
		if(ctx.IsGet()) {
			path := string(ctx.Path());
			log.Printf("request path: %s", path);
			pathElements := strings.Split(path, "/")[1:];
			elementCount := len(pathElements);

			switch {
			case elementCount == 2 && pathElements[0] == "date":
				datestring := pathElements[1];
				log.Printf("getting data for date %s", datestring);
				RenderDatePage(ctx, datestring, c);
			case elementCount == 1 && pathElements[0] == "":
				log.Printf("getting availability page");
				RenderAvailabilityPage(ctx, c);
			default:
				log.Printf("invalid path requested:  %s", path);
				ctx.Error("requested path is invalid", fasthttp.StatusNotFound);
			}
		}
	}

	log.Printf("Starting the exercise server on port %s", port);
	err := fasthttp.ListenAndServe(port, server);
	exitMessage := "exercise server has stopped";
	if err != nil {
		exitMessage = fmt.Sprintf("%s with error: %s", exitMessage, err);
	}
	log.Fatalf("%s", exitMessage);
}