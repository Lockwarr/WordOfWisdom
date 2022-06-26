package cmd

import (
	"context"
	"sitemapGenerator/WordOfWisdom/internal/client"
)

func main() {
	err := client.Run(context.Background(), ":8080")
	if err != nil {
		panic(err)
	}
}
