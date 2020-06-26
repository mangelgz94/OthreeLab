package main

func main() {
	application := App{}

	application.Initialize()
	application.Run(":8081")
}
