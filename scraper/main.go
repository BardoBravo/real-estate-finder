package cloudrun

import (
	"context"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gocolly/colly"
	"google.golang.org/api/option"
)

var DBAccess = startDBService()

type Platform interface {
	crawl(config Config, writer Exporter) *colly.Collector
}

type Item struct {
	//TODO: version 1.1 -> add new fields
	Title            string
	Location         string
	HasExactLocation bool
	//price            int
	Price string
	//livingSpace float64
	LivingSpace string
	Rooms       float64
	Url         string
	ScrapedAt   time.Time
}

type Config struct {
	dataDir          string
	platforms        []Platform
	storage          Storage
	isDefined        bool
	collectorOptions []colly.CollectorOption
}

type DBCall struct {
	DBContext context.Context
	DBClient  *firestore.Client
}

func (record Item) csvRow() []string {
	return []string{
		record.Title,
		record.Location,
		strconv.FormatBool(record.HasExactLocation),
		//strconv.Itoa(record.price),
		record.Price,
		//strconv.FormatFloat(record.livingSpace, 'f', -1, 64),
		record.LivingSpace,
		strconv.FormatFloat(record.Rooms, 'f', -1, 64),
		record.Url,
		record.ScrapedAt.Format(time.RFC3339),
	}
}

func readConfig(params url.Values) Config {
	available := map[string]Platform{
		//"ebay_kleinanzeigen": EBayKleinanzeigen{},
		"jus_sl": Jus_sl{},
		//"rodrigues": Rodrigues{},
		//"nestpick":           Nestpick{},
	}

	platforms := make([]Platform, 0)
	for name := range available {
		platforms = append(platforms, available[name])
	}

	cache := params.Get("cache") == "1"
	var collectorOptions []colly.CollectorOption
	if cache {
		collectorOptions = append(collectorOptions, colly.CacheDir("cache"))
	}
	platform := params.Get("platform")
	if platform != "" {
		platforms = []Platform{available[platform]}
	}

	bucket, isDefined := os.LookupEnv("GCLOUD_BUCKET")
	if !isDefined && os.Getenv("PORT") == "8080" {
		log.Fatalln("GCLOUD_BUCKET must be defined")
	}
	date := time.Now().UTC().Format(time.RFC3339)
	storage := GCloudStorage{
		bucket:          bucket,
		destinationPath: date + "/",
	}

	return Config{
		dataDir:          "results",
		platforms:        platforms,
		storage:          storage,
		isDefined:        isDefined,
		collectorOptions: collectorOptions,
	}
}

func startDBService() (dbAccess DBCall) {

	var ctx = context.Background()
	_, err := os.Stat("deployment/credentials.json")

	if os.IsNotExist(err) {
		var conf = &firebase.Config{ProjectID: "find-new-rent"}
		var app, errorApp = firebase.NewApp(ctx, conf)
		if errorApp != nil {
			log.Fatalf("app./firestore: %v", errorApp)
		}
		var client, errFirestore = app.Firestore(ctx)
		if errFirestore != nil {
			log.Fatalf("firebase.NewClient: %v", errorApp)
		}
		//defer client.Close()
		return DBCall{
			DBClient:  client,
			DBContext: ctx,
		}

	} else {
		var sa = option.WithCredentialsFile("deployment/credentials.json")
		var app, errorApp = firebase.NewApp(ctx, nil, sa)
		if errorApp != nil {
			log.Fatalf("firebase.NewApp: %v", errorApp)
		}
		var client, errFirestore = app.Firestore(ctx)
		if errFirestore != nil {
			log.Fatalf("firebase.NewClient: %v", errorApp)
		}
		//defer client.Close()
		return DBCall{
			DBClient:  client,
			DBContext: ctx,
		}

	}

}

func Run(params url.Values) string {
	config := readConfig(params)
	for _, platform := range config.platforms {
		fileName := strings.Split(reflect.TypeOf(platform).String(), ".")[1]
		sender := EmailSender{scraperName: fileName}
		fileName = path.Join(config.dataDir, fileName+".csv")
		exporter := CSVExporter{fileName: fileName}
		exporter.run(config, platform.crawl)
		//TODO: version 1.3 -> just store the new items
		config.storage.write(fileName)
		sender.writeAndSend()
	}

	return "it works"
}
