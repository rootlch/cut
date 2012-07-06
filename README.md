cut is a very simple utility that can help your perform different kinds of tasks. cut implements the io.Writer interface so you can put data into it using the fmt package.


For example:
You can break a file into lines and process each line individually.

file, err := os.Open("file.go")
if err != nil {
	log.Fatal(err)
}
c := cut.New()
os.Copy(c, file)

for line := range c.Between("", "\n") {
	fmt.Println(line)
}


You can also use it to scrape all links from a webpage.

resp, err := http.Get("http://example.com/")
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
c := cut.New()
os.Copy(c, resp.Body)

for line := range c.Between("<a ", "a>") {
	fmt.Println(line)
}
