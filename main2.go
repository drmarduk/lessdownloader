package less

import (
	"flag"
	"fmt"
)

// Global Variables parsed from args
// may be improved to be local in main and passed to needed functions
// well first try.
var strAction string
var strID string
var intMax int


func main() {
	// gather commandline infos
	var help = flag.Bool("help", false, "Help message for shit")
	flag.StringVar(&strAction, "action", "download", "The Action to be executed: download, ...cont.")
	flag.StringVar(&strID, "id", "empty", "This is the 'ID' of either a Picture, Gallery or Video.")
	flag.IntVar(&intMax, "max", 0, "Maximum entries to be {actioned}. Default 0 -> infinite.")



	flag.Parse()

	fmt.Println("help has value", *help)

}