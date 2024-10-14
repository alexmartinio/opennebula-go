package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/image"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/image/keys"
	"github.com/joho/godotenv"
)

func getVMs(controller *goca.Controller) {
	// Get short informations of the VMs
	vms, err := controller.VMs().Info()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of VMs: %s\n\n", strconv.Itoa(len(vms.VMs)))

	file, err := os.Create("opennebula-vms.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	line := "ID, VM Name, Username, State, User Details\n"
	fmt.Print(line)
	fmt.Fprint(file, line)

	for i := 0; i < len(vms.VMs); i++ {
		// This Info method, per VM instance, give us detailed informations on the instance
		// Check xsd files to see the difference
		vm, err := controller.VM(vms.VMs[i].ID).Info(false)
		if err != nil {
			log.Fatal(err)
		}

		// Do some others stuffs on vm
		name := ""
		state, _, _ := vm.State()

		line := fmt.Sprintf("%+v, %+v, %+v, %+v, %+v\n", vm.ID, vm.Name, vm.UName, state, name)
		fmt.Print(line)
		fmt.Fprint(file, line)
	}
}

func getImages(controller *goca.Controller) {
	// Get short informations of the Images
	images, err := controller.Images().Info()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of Images: %s\n\n", strconv.Itoa(len(images.Images)))

	file, err := os.Create("opennebula-images.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	line := "ID, VM Name, Username, State, User Details\n"
	fmt.Print(line)
	fmt.Fprint(file, line)

	for i := 0; i < len(images.Images); i++ {
		// This Info method, per Image instance, give us detailed informations on the image
		// Check xsd files to see the difference
		image, err := controller.Image(images.Images[i].ID).Info(false)
		if err != nil {
			log.Fatal(err)
		}

		// Do some others stuffs on vm
		name := ""
		state, _ := image.State()

		line := fmt.Sprintf("%+v, %+v, %+v, %+v, %+v\n", image.ID, image.Name, image.UName, state, name)
		fmt.Print(line)
		fmt.Fprint(file, line)
	}
}

func getMarketApps(controller *goca.Controller) {
	// Get short informations of the Marketplace Apps
	apps, err := controller.MarketPlaceApps().Info()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of Market Place Apps: %s\n\n", strconv.Itoa(len(apps.MarketPlaceApps)))

	app_id, err := controller.MarketPlaceApps().ByName("Alpine Linux 3.8")
	if err != nil {
		log.Fatal(err)
	}

	tpl := image.NewTemplate()
	tpl.Add(keys.Name, "test-image3")
	tpl.Add(keys.Size, "4096")
	tpl.Add("FROM_APP", app_id)
	tpl.Add(keys.Path, "https://marketplace.opennebula.io//appliance/f4cc1890-f430-013c-b669-7875a4a4f528/download/0")

	fmt.Print(tpl.String())

	fmt.Print(app_id)

	app, _ := controller.MarketPlaceApp(app_id).Info(true)
	fmt.Print(app)

	image, err := controller.Images().Create(tpl.String(), 100)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(image)
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := goca.NewDefaultClient(
		goca.NewConfig(
			os.Getenv("ONE_USERNAME"),
			os.Getenv("ONE_TOKEN"),
			os.Getenv("ONE_ENDPOINT"),
		),
	)
	controller := goca.NewController(client)
	getVMs(controller)
	getImages(controller)
	getMarketApps(controller)
}
