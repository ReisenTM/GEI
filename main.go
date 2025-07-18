package main

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
import (
	"GEI/gei"
	"fmt"
	"net/http"
)

func main() {
	r := gei.New()
	r.GET("/test", func(c *gei.Context) {
		c.JSON(http.StatusOK, gei.H{
			"hello": "world",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
